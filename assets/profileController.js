// new dependency: ngResource is included just above
// var myapp = new angular.module("myapp", ["ngResource"]);

// inject the $resource dependency here
myapp.controller("ProfileController", ["$scope", "$http", function($scope, $http){
    // I designed the backend to play nicely with angularjs so this is all the
    // setup we need to do all of the ususal operations.


    var oldSidebarColor;

        // if ( modalLoading.isActive() ) {
        //     modalLoading.hide();
        // } else {
        //     modalLoading.show();
        // }

    $scope.getProfile = function() {
        setLoadingScreen(true);
        $http.get("/api/profile").then(function (response) {
            $scope.profile = response.data;
            setLoadingScreen(false);
            oldSidebarColor = $scope.profile.SidebarColor;
        });
    }

    $scope.updateProfile = function() {


        setLoadingScreen(true);
        $http.post('/api/profile', $scope.profile
        ).success(function(data, status, headers, config) {

            if (oldSidebarColor != $scope.profile.SidebarColor) {
                location.reload();
            } else {
                $scope.getProfile();
                setLoadingScreen(false);
            }

        }).error(function(data, status) {
            UIkit.notify("Error! "+status+"<br/><small>"+data+"</small>");
        });
    }
    $scope.updateAdvancedSettings = function() {
        $scope.updateProfile();

    }

    $scope.updatePassword = function() {

        if ($scope.password == $scope.passwordCheck) {
            $http.post('/api/profile/password', { password: $scope.password }
            ).success(function(data, status, headers, config) {
                UIkit.modal.alert("Done!");
            }).error(function(data, status) {
                UIkit.modal.alert("Error! Please tell <a href='http://telegram.me/mdeheij/'>@mdeheij</a> about this!<br/><small>"+data+"</small>");
            });
        } else {
            UIkit.modal.alert("Error! Passwords do not match");
            $scope.password = "";
            $scope.passwordCheck = "";
        }
    }

    $scope.getProfile();
    // var Invoice = $resource("/api/profile/:id", {id: '@id'}, {});
    //
    // $scope.selected = null;
    // $scope.selectedFinishedInvoice = false;
    //
    // $scope.list = function(idx){
    //     // Notice calls to Invoice are often given callbacks.
    //     Invoice.query(function(data){
    //         $scope.invoices = data;
    //     }, function(error){
    //         alert(error.data);
    //     });
    // };
    //
    // $scope.list();
    //
    // $scope.get = function(idx){
    //     console.log("[get] "+idx);
    //     // Passing parameters to Invoice calls will become arguments if
    //     // we haven't defined it as part of the path (we did with id)
    //     Invoice.get({id: $scope.invoices[idx].id}, function(data){
    //         $scope.selected = data;
    //         $scope.selected.idx = idx;
    //     });
    // };
    //
    // $scope.update = function(invoice) {
    //         invoice.$save().then(function() {
    //              console.log("Daarna");
    //              $scope.list();
    //             }
    //          );
    // };

}]);
