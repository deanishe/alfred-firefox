/* global browser */

/**
 * Name of native application according to application manifest.
 * @var {string} appName
 */
const appName = 'net.deanishe.alfred.firefox';

const iconConnected = 'icons/bowler.svg';
const iconDisconnected = 'icons/bowler-red.svg';

/**
 * Tab object.
 * @param {tabs.Tab} tab - Native tab object to create Tab from.
 * @return {Object} - API Tab object.
 */
const Tab = tab => {
  let obj = {};

  tab = tab || {};

  obj.id         = tab.id          || 0;
  obj.windowId   = tab.windowId    || 0;
  obj.index      = tab.index       || 0;
  obj.title      = tab.title       || '';
  obj.url        = new URL(tab.url || '');
  // obj.favicon = tab.favIconUrl  || '';
  obj.active     = tab.active      || false;

  obj.toString = function() {
    return `#${this.id} (${this.windowId}x${this.index}) "${this.title}" - ${this.url}`;
  };

  return obj;
};

/**
 * Bookmark object.
 * @param {bookmarks.BookmarkTreeNode} bm - Native object to create Bookmark from.
 * @return {Object} - API Bookmark object.
 */
const Bookmark = bm => {
  let obj = {};
  bm = bm || {};

  obj.id       = bm.id       || 0;
  obj.index    = bm.index    || 0;
  obj.title    = bm.title    || '';
  obj.parentId = bm.parentId || 0;
  obj.type     = bm.type     || '';
  obj.url      = bm.url      || '';

  obj.toString = function() {
    return `#${this.id} "${this.title}" - ${this.url}`;
  };

  return obj;
};

/**
 * HistoryEntry object.
 * @param {history.HistoryItem} hi - Native object to create HistoryEntry from.
 * @return {Object} - API History object.
 */
const HistoryEntry = hi => {
  let obj = {};
  hi = hi || {};

  obj.id    = hi.id    || 0;
  obj.url   = hi.url   || '';
  obj.title = hi.title || hi.url;

  obj.toString = function() {
    return `#${this.id} "${this.title}" - ${this.url}`;
  };

  return obj;
};

/**
 * Download object.
 * @param {downloads.DownloadItem} di - Native object to create Download from.
 * @return {Object} - API Download object.
 */
const Download = di => {
  let obj = {};
  di = di || {};

  obj.id     = di.id       || 0;
  obj.path   = di.filename || '';
  obj.size   = di.fileSize || 0;
  obj.url    = di.url      || '';
  obj.mime   = di.mime     || '';
  obj.exists = di.exists   || false;
  obj.error  = di.error    || '';

  obj.toString = function() {
    return `#${this.id} "${this.path}" - ${this.url}`;
  };

  return obj;
};

/**
 * Extension application object.
 * @constructor
 */
const Background = function() {
  const self = this;

  self.port = null,
    self.nativePort = null,
    self.connected = false;

  self.onConnected = port => {
    self.port = port;
    port.onMessage.addListener(self.receive);
    console.debug('connected to popup');
  };

  self.send = msg => {
    self.port.postMessage(msg);
    console.debug('sent message', msg);
  };

  self.receive = msg => {
    console.debug('received message', msg);
    if ('command' in msg) {
      switch (msg.command) {
        case 'status':
          self.send({ status: self.connected ? 'connected' : 'disconnected' });
          return;
        case 'reconnect':
          console.debug('reconnecting to native app ...');
          self.connectNative();
          return;
        case 'reload':
          console.debug('reloading extension ...');
          browser.runtime.reload();
          return;
      }
    }
  };

  self.connectNative = () => {
    self.connected = false;

    let listener = payload => {
      if (!self.connected) {
        self.connected = true;
        // self.nativePort.onDisconnect.removeListener(self.connectNativeFailed);
        self.onConnectedNative();
      }
      self.receiveNative(payload);
    };

    self.nativePort = browser.runtime.connectNative(appName);
    self.nativePort.onMessage.addListener(listener);
    self.nativePort.onDisconnect.addListener(self.connectNativeFailed);
  };

  /**
   * Callback for connection failure.
   * Logs an error message to the console.
   */
  self.connectNativeFailed = port => {
    let msg = '';
    if (port.error) {
      msg = port.error.message;
    }
    self.connected = false;
    console.error(`native client connection failed: ${msg}`);
    browser.browserAction.setIcon({ path: iconDisconnected });
  };

  /**
   * Callback for successful connection to native application.
   * Logs a message to the console.
   */
  self.onConnectedNative = () => {
    console.log('connected to native client');
    browser.browserAction.setIcon({ path: iconConnected });
  };

  /**
   * Handle commands from native application.
   * @param {Object} msg - Data from native application.
   * @param {string} msg.id - Command/response ID.
   * @param {Object} msg.params - Arguments to command.
   */
  self.receiveNative = msg => {
    console.log(`received:`, msg);
    let p = null;
    if ('command' in msg) {
      switch (msg.command) {
        case 'ping':
          p = self.ping();
          break;
        // case 'all-windows':
        //   p = self.allWindows();
        //   break;
        // case 'current-window':
        //   p = self.currentWindow();
        //   break;
        case 'all-tabs':
          p = self.allTabs();
          break;
        // DEPRECATED - replaced by self.tab(); unused by newer
        // versions 0.2.0+ of workflow
        // Remove from future versions
        case 'current-tab':
          p = self.tab(0);
          break;
        case 'tab':
          p = self.tab(msg.params);
          break;
        case 'all-bookmarks':
          p = self.allBookmarks();
          break;
        case 'search-bookmarks':
          p = self.searchBookmarks(msg.params);
          break;
        case 'search-history':
          p = self.searchHistory(msg.params);
          break;
        case 'search-downloads':
          p = self.searchDownloads(msg.params);
          break;
        case 'activate-tab':
          p = self.activateTab(msg.params);
          break;
        case 'close-tabs-left':
          p = self.closeTabsLeft(msg.params);
          break;
        case 'close-tabs-right':
          p = self.closeTabsRight(msg.params);
          break;
        case 'close-tabs-other':
          p = self.closeTabsOther(msg.params);
          break;
        case 'execute-js':
          p = self.executeJS(msg.params);
          break;
        case 'run-bookmarklet':
          p = self.runBookmarklet(msg.params);
          break;
        case 'open-incognito':
          p = self.openIncognito(msg.params);
          break;
        default:
          console.error(`unknown command: ${msg.command}`);
          self.sendError(msg.id, 'unknown command');
          return;
      }
      p.then(payload => {
        self.sendNative({ id: msg.id, payload: payload });
      }).catch(err => {
        self.sendError(msg.id, err.message);
      });
    } else {
      self.sendError(msg.id, 'no command given');
    }
  };

  /**
   * Send response to native application.
   * @param {Object} msg - Data to send to native application.
   * @param {string} msg.id - Command/response ID.
   * @param {string|bool|Object} msg.payload - Actual response data.
   * @param {string} msg.error - Error message if command failed.
   */
  self.sendNative = msg => {
    if (self.nativePort) {
      self.nativePort.postMessage(msg);
      console.log(`sent:`, msg);
    }
  };

  /**
   * Send error respones to native application.
   * @param {string} id - Command/response ID.
   * @param {string} msg - Error message.
   */
  self.sendError = (id, msg) => {
    self.sendNative({ id: id, error: msg });
  };

  /**
   * Handle "ping" command.
   * @return {Promise} - Resolves to string "pong".
   */
  self.ping = () => {
    return new Promise(resolve => {
      resolve('pong');
    });
  };

  /**
   * Handle "all-tabs" command.
   * @return {Promise} - Resolves to array of Tab objects for all tabs sorted
   * by most recently used.
   */
  self.allTabs = () => {
    return browser.tabs.query({}).then(tabs => tabs
      .sort((a, b) => (b?.lastAccessed ?? 0) - (a?.lastAccessed ?? 0))
      .map(t => Tab(t))
    );
  };

  /**
   * Handle "activate-tab" command.
   * @param {number} id - ID of tab to activate.
   */
  self.activateTab = id => {
    return browser.tabs
      .update(id, { active: true })
      .then(() => {
        return browser.tabs.get(id);
      })
      .then(tab => {
        return browser.windows.update(tab.windowId, { focused: true });
      });
  };

  /**
   * Handle "current-tab" command.
   * @return {Promise} - Resolves to Tab for current tab.
   * Throws an error if there is no current tab.
   */
  // self.currentTab = () => {
  //   return self.activeTab(null).then(t => {
  //     if (!t) throw 'no current tab';
  //     let tab = Tab(t);
  //     console.log(`[current-tab] ${tab}`);
  //     return tab;
  //   });
  // };

  /**
   * Handle "tab" command.
   * @param {number} tabId - ID of tab to return.
   * @return {Promise} - Resolves to Tab for current tab.
   * Throws an error if there is no current tab.
   */
  self.tab = tabId => {
    if (!tabId) {
      return self.activeTab(null).then(t => {
        if (!t) throw 'no current tab';
        let tab = Tab(t);
        console.log(`[current-tab] ${tab}`);
        return tab;
      });
    }

    return browser.tabs
      .get(tabId)
      .then(t => {
        return Tab(t);
      })
  };

  /**
   * Handle "all-bookmarks" command.
   * @return {Promise} - Resolves to array of Bookmark objects for all bookmarks
   * and folders.
   */
  self.allBookmarks = () => {
    let bookmarks = [];
    let addBookmarks = node => {
      if (node.url) bookmarks.push(Bookmark(node));
      if (node.children) node.children.map(n => addBookmarks(n));
    };

    return browser.bookmarks.getTree().then(root => {
      addBookmarks(root[0]);
      return bookmarks;
    });
  };

  /**
   * Handle "search-bookmarks" command.
   * @param {string} query - Search query.
   * @return {Promies} - Resolves to array of Bookmark objects matching query.
   */
  self.searchBookmarks = query => {
    return browser.bookmarks.search(query).then(nodes => {
      let bookmarks = nodes.filter(n => n.url).map(n => Bookmark(n));
      console.debug(`${bookmarks.length} bookmark(s) for "${query}"`);
      return bookmarks;
    });
  };

  /**
   * Handle "search-history" command.
   * @param {string} query - Search query.
   * @return {Promies} - Resolves to array of History objects matching query.
   */
  self.searchHistory = query => {
    return browser.history.search({ text: query, startTime: 0 }).then(items => {
      let history = items.filter(it => it.url).map(it => HistoryEntry(it));
      console.debug(`${history.length} history item(s) for "${query}"`);
      return history;
    });
  };

  /**
   * Handle "search-downloads" command.
   * @param {string} query - Search query.
   * @return {Promise} - Resolves to array of Download objects matching query.
   */
  self.searchDownloads = query => {
    return browser.downloads
      .search({
        query: [query],
        exists: true,
      })
      .then(items => {
        console.debug(`${items.length} download(s) for "${query}"`);
        return items.map(it => Download(it));
      });
  };

  /**
   * Handle "close-tabs-left" command.
   * @param {number} tabId - ID of tab whose neighbours to close.
   * @return {Promise} - Result of browser.tabs.remove()
   */
  self.closeTabsLeft = tabId => {
    console.debug(`closing tabs to left of tab #${tabId} ...`);
    let activeTab = null;
    return browser.tabs
      .get(tabId)
      .then(tab => {
        if (!tab) throw 'no current tab';
        activeTab = tab;
        return browser.tabs.query({ windowId: tab.windowId });
      })
      .then(tabs => {
        let ids = tabs.filter(t => t.index < activeTab.index).map(t => t.id);
        return browser.tabs.remove(ids);
      });
  };

  /**
   * Handle "close-tabs-right" command.
   * @param {number} tabId - ID of tabs whose neighbours to close.
   * @return {Promise} - Result of browser.tabs.remove()
   */
  self.closeTabsRight = tabId => {
    console.debug(`closing tabs to right of tab #${tabId} ...`);
    let activeTab = null;
    return browser.tabs
      .get(tabId)
      .then(tab => {
        if (!tab) throw 'no current tab';
        activeTab = tab;
        return browser.tabs.query({ windowId: tab.windowId });
      })
      .then(tabs => {
        let ids = tabs.filter(t => t.index > activeTab.index).map(t => t.id);
        return browser.tabs.remove(ids);
      });
  };

  /**
   * Handle "close-tabs-other" command.
   * @param {number} tabId - ID of window to close tabs in.
   * @return {Promise} - Result of browser.tabs.remove()
   */
  self.closeTabsOther = tabId => {
    console.debug(`closing other tabs in window of tab #${tabId} ...`);
    let activeTab = null;
    return browser.tabs
      .get(tabId)
      .then(tab => {
        activeTab = tab;
        return browser.tabs.query({ windowId: tab.windowId });
      })
      .then(tabs => {
        let ids = tabs.filter(t => t.id !== activeTab.id).map(t => t.id);
        return browser.tabs.remove(ids);
      });
  };

  /** Handle "execute-js" command. */
  // self.executeJS = js => {
  //   return browser.tabs.executeScript({ code: js }).then(results => {
  //     console.debug(`js=${js}, results=`, results);
  //   });
  // };
  /**
   * Handle "execute-js" command.
   * @param {Object} params - Tab and bookmarklet IDs.
   * @param {number} params.tabId - ID of tab to execute JS in.
   * If tabId is 0, JS is executed in the active tab.
   * @param {string} params.js - JavaScript to execute.
   */
  self.executeJS = params => {
    console.debug(`execute-js`, params);
    var p;
    if (params.tabId) {
      p = browser.tabs.executeScript(params.tabId, { code: params.js });
    } else {
      p = browser.tabs.executeScript({ code: params.js });
    }
    return p.then(result => {
      return JSON.stringify(result);
    });
  };

  /**
   * Handle "run-bookmarklet" command.
   * @param {Object} params - Tab and bookmarklet IDs.
   * @param {number} params.tabId - ID of tab to execute bookmarklet in.
   * If tabId is 0, bookmarklet is executed in the active tab.
   * @param {string} params.bookmarkId - ID of bookmarklet to execute.
   */
  self.runBookmarklet = params => {
    console.debug(`run-bookmarklet`, params);
    return browser.bookmarks.get(params.bookmarkId).then(bookmarks => {
      if (!bookmarks.length) throw 'bookmark not found';
      let bm = bookmarks[0];
      if (!bm.url.startsWith('javascript:')) throw 'not a bookmarklet';
      let js = decodeURI(bm.url.slice(11));
      if (params.tabId) browser.tabs.executeScript(params.tabId, { code: js });
      else browser.tabs.executeScript({ code: js });
    });
  };

  /**
   * Handle "open-incognito" command.
   * @param {string} url - URL to open in a new Incognito window.
   * @return {Promise} - Promise that resolves to null.
   */
  self.openIncognito = url => {
    console.debug(`open-incognito ${url}`);
    return browser.windows.create({ incognito: true, url: url });
  };

  /**
   * Return active tab.
   * @param {number} winId - ID of window to get active tab of.
   * If 0 or null, current window is used.
   * @return {Promise} - Promise resolves to null or a tabs.Tab.
   */
  self.activeTab = winId => {
    winId = winId || browser.windows.WINDOW_ID_CURRENT;
    return browser.tabs
      .query({
        active: true,
        windowId: winId,
      })
      .then(tabs => {
        if (tabs.length) return tabs[0];
        return null;
      });
  };

  browser.runtime.onConnect.addListener(self.onConnected);
  self.connectNative();
  console.log(`started`);
};

browser.browserAction.setIcon({ path: iconDisconnected });
new Background();
