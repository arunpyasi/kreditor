// new dependency: ngResource is included just above
// var myapp = new angular.module("myapp", ["ngResource"]);

// inject the $resource dependency here
myapp.controller("InvoiceController", ["$scope", "$resource", function($scope, $resource){
    // I designed the backend to play nicely with angularjs so this is all the
    // setup we need to do all of the ususal operations.
    var Invoice = $resource("/api/invoice/:id", {id: '@id'}, {});

    var addModal = UIkit.modal("#newInvoice");

    $scope.selected = null;
    $scope.selectedFinishedInvoice = false;

    $scope.list = function(idx){
        // Notice calls to Invoice are often given callbacks.
        Invoice.query(function(data){
            $scope.invoices = data;
            // if(idx != undefined) {
            //     $scope.selected = $scope.invoices[idx];
            //     $scope.selected.idx = idx;
            // }
        }, function(error){
            alert(error.data);
        });
    };

    $scope.list();

    $scope.get = function(idx){
        console.log("[get] "+idx);
        // Passing parameters to Invoice calls will become arguments if
        // we haven't defined it as part of the path (we did with id)
        Invoice.get({id: $scope.invoices[idx].id}, function(data){
            $scope.selected = data;
            $scope.selected.idx = idx;
        });
    };


    $scope.openAddModel = function(idx) {

        addModal.show();
        $scope.modalTitle = "Add new invoice";
        $scope.modalAction = "add";
        $scope.selectedObject = new Invoice();
        $scope.selectedObject.Currency = "€";

    }


    $scope.add = function() {
        var newInvoice = new Invoice();
        addModal.show();



            if ($scope.selectedObject.Debtor != undefined) {
              newInvoice.Debtor = $scope.selectedObject.Debtor;
              newInvoice.IncludeLinks = $scope.selectedObject.IncludeLinks;
              newInvoice.$save().then(function() {
                   console.log("Moving on..");
                   $scope.list();
                   addModal.hide();
                  }
               );
            } else {
              UIkit.modal.alert("You need to specify a debtor!");
            }



    };

    $scope.update = function(invoice) {
        UIkit.modal.confirm('Are you sure you want to regenerate this link?', function(){
            console.log("Commencing update();...");
            console.log("Using invoice:");
            console.log(invoice);

            console.log("Saving..");

            invoice.$save().then(function() {
                 console.log("Daarna");
                 $scope.list();
                }
             );
        });
    };
    function OpenInNewTab(url) {
        var win = window.open(url, '_blank');
        win.focus();
    }

    $scope.openInvoice =  function(link) {
        OpenInNewTab("/i/"+link);
    }

    $scope.remove = function(object){
        UIkit.modal.confirm('Are you sure?', function(){
            addModal.hide();
            object.$delete().then(function() {
                 console.log("Deleting..");
                 $scope.list();
                }
             );
        });
    };
}]);
