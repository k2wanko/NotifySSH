(function(){


  document.addEventListener("DOMContentLoaded", function(event) {

    //
    // Polymer toggle
    //
    var notify_toggle = document.querySelector("paper-toggle-button");
    notify_toggle.addEventListener('change', function () {
      if (this.checked) {
        //
        // ServiceWorker Reg
        //
        navigator.serviceWorker.register('service-worker.js').then(function(sw){

          return sw.pushManager.subscribe().then(function(sub){
            console.log(sub);
          });

        }).catch(function(err){
          console.error(err);
        });

        document.getElementById("notify_box").style.display="block";

      }else{
        document.getElementById("notify_box").style.display="none";
      }

    });

  });

})();
