
// inject the $resource dependency here
myapp.controller("MainCtl", ["$scope", "$resource", "$http", function($scope, $resource, $http){
    // I designed the backend to play nicely with angularjs so this is all the
    // setup we need to do all of the ususal operations.
    var Debt = $resource("/api/debts/:id", {id: '@id'}, {});

    //

    var addModal = UIkit.modal("#modal6");

    $scope.selected = null;
    $scope.selectedFinishedDebt = false;

    $scope.list = function(idx){


        // Notice calls to Debt are often given callbacks.
        Debt.query(function(data){
            $scope.debts = data;
            // if(idx != undefined) {
            //     $scope.selected = $scope.debts[idx];
            //     $scope.selected.idx = idx;
            // }
        }, function(error){
            alert(error.data);
        });

        $http.get("/api/debtlinks").then(function (response) {
            $scope.debtlinks = response.data;
        });

    };

    $scope.negativeValue=function(myValue){
        if(parseInt(myValue) < 0){
                var css = { 'color':'red' };
                return css;
        }
    }
    $scope.arrowIcon=function(myValue){
        var num = parseInt(myValue);
        if(num < 0){
            var icon = 'fa-arrow-up';
            return icon;
        } else {
            var icon = 'fa-arrow-down';
            return icon;
        }
    }

    $scope.list();

    $scope.get = function(idx){
        console.log("[get] "+idx);
        // Passing parameters to Debt calls will become arguments if
        // we haven't defined it as part of the path (we did with id)
        Debt.get({id: $scope.debts[idx].id}, function(data){
            $scope.selected = data;
            $scope.selected.idx = idx;
        });
    };

    $scope.debtDirection = function() {
        // var str = $scope.selectedObject.Amount;
        // if (str.indexOf('-') !== -1) {
        //     console.log("hij was negatief.");
        // }

        $scope.selectedObject.Amount = $scope.selectedObject.Amount*-1;


        console.log("Debt direction changed!");
    }

    $scope.openAddModel = function(idx) {

        addModal.show();
        $scope.modalTitle = "Add new debt";
        $scope.modalAction = "add";
        $scope.selectedObject = new Debt();
        $scope.selectedObject.Currency = "€";

    }

    $scope.openUpdateModal = function(idx) {
        //Use Angular's copy method here to ensure the original object does not get updated if the user cancels
        $scope.selectedObject = angular.copy(idx);

        $scope.modalTitle = "Edit this debt";
        $scope.modalAction = "edit";
        addModal.show();

    }

    $scope.add = function() {
        // I was lazy with the user input.
        //var modal = UIkit.modal("#modal6");
        //modal.show();

        // var title = prompt("Enter the debt's title.");
        // if(title == null){
        //     return;
        // }
        // var author = prompt("Enter the debt's author.");
        // if(author == null){
        //     return;
        // }
        // Creating a blank debt object means you can still $save
        var newDebt = new Debt();
        newDebt.Amount = $scope.selectedObject.Amount;
        newDebt.Description = $scope.selectedObject.Description;
        newDebt.Debtor = $scope.selectedObject.Debtor;
        newDebt.Paid = $scope.selectedObject.Paid;
        //newDebt.title = title;
        //newDebt.author = author;


        newDebt.$save().then(function() {
             console.log("Daarna");
             $scope.list();
             addModal.hide();
            }
         );
    };

    $scope.update = function(debt) {


        //debt.title = title;
//
/*
var newDebt = new Debt();

*/
    console.log("Commencing update();...");
    console.log("Using debt:");
    console.log(debt);
        debt.Amount = $scope.selectedObject.Amount;
        debt.Description = $scope.selectedObject.Description;
        debt.Debtor = $scope.selectedObject.Debtor;
        debt.Paid = $scope.selectedObject.Paid;
        // Noticed I never created a new Debt()?

        console.log("Saving..");
        addModal.hide();
        //newDebt.author = author;
        //TODO: Check if this is bad practice whenever my plane lands and I got internet again
        debt.$save().then(function() {
             console.log("Daarna");
             $scope.list();
            }
         );

        //debt.$save()
        //$scope.list();
        //$scope.list(debt.$save());
    };
    //
    // $scope.remove = function(idx){
    //     $scope.debts[idx].$delete();
    //     $scope.selected = null;
    //     $scope.list();
    // };

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
