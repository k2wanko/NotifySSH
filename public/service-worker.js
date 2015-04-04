(function(){

  self.addEventListener('push', function(e){

    console.log('on push', e);
    
    var title = "NotifySSH";
    var body = "Hey";
    var icon = "/images/icon.png";
    var tag = "notifyssh-login";
    
    e.waitUntil(self.registration.showNotification(title, {
      body: body,
      icon: icon,
      tag: tag
    }));
  });

  self.addEventListener('notificationclick', function(e){
    console.log('on notificationclick', e);

    e.notification.close();
    e.waitUntil(
      clients.openWindow('/')
    );
  });
  
})();
