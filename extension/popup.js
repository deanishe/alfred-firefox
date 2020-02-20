/* global browser */
const Popup = function() {
  const self = this;

  self.status = null;

  self.connected = () => {
    return self.status === true;
  };

  self.receive = msg => {
    if ('status' in msg) {
      self.status = msg.status;
      console.debug(`[popup] status=${msg.status}`);
    }
  };
  self.send = msg => {
    self.port.postMessage(msg);
    // console.debug(`[popup] sent message`, msg);
  };

  self.reconnect = () => {
    self.send({ command: 'reconnect' });
  };

  self.reload = () => {
    self.send({ command: 'reload' });
  };

  self.getStatus = () => {
    self.send({ command: 'status' });
  };

  self.port = browser.runtime.connect({ name: 'alfred-firefox' });
  self.port.onMessage.addListener(self.receive);
  console.debug(`[popup] started`);
  self.getStatus();
};

(function() {
  const popup = new Popup();
  const status = document.querySelector('#status');

  // document
  //   .querySelector('#reloadButton')
  //   .addEventListener('click', popup.reload);

  setInterval(function() {
    popup.getStatus();
    switch (popup.status) {
      case 'connected':
        status.classList.remove('disconnected');
        status.classList.add('connected');
        break;
      case 'disconnected':
        status.classList.remove('connected');
        status.classList.add('disconnected');
        popup.reconnect();
        break;
    }
  }, 500);
})();
