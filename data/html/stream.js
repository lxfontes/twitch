Vue.config.devtools = true;
Vue.config.debug = true;

var data = {};
function sharedData(merge) {
  for (var k in merge) {
    data[k] = merge[k];
  }

  return function() {
    return data;
  };
};

var Twitch = window.Twitch = {
  register: function(component, opts) {
    return Vue.component(component, opts);
  },
};

Twitch.register('blank', {
  template: '<p>ajkasdjfkajsdfkj</p>',
});

var vm = new Vue({
  el: '#app',
  template: '<div id="sandbox"><component :is="currentView"></component></div>',
  data: {
    ws: null,
    currentView: Vue.component('blank'),
  },
  created: function() {
    this.setupWebsocket();
  },
  methods: {
    onMessage: function(evtRaw) {
      var ev = JSON.parse(evtRaw.data);
      sharedData(ev.args);
      this.currentView = Vue.component(ev.name);

      var self = this;
      setTimeout(function() {
        self.currentView = Vue.component('blank');
      }, 5000);
    },
    setupWebsocket: function() {
      var address = 'ws://' + location.host + '/ws';
      console.log("Connecting to " + address);

      var ws = new WebSocket(address);

      ws.onopen = function() {
        console.log("ws connected");
      },

      ws.onmessage = this.onMessage;

      ws.onclose = function() {
        setTimeout(this.setupWebsocket, 5000);
      }.bind(this);
    }
  }
});
