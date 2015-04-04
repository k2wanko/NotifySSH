(function(){


  document.addEventListener("DOMContentLoaded", function(event) {
  
    navigator.serviceWorker.register('service-worker.js').then(function(){
      
      return navigator.serviceWorker.ready.then(function(sw){
        
        sw.pushManager.subscribe().then(function(sub){
          console.log(sub);
        });
      });
      
    }).catch(function(err){
      console.error(err);
    });
    
  });
})();






