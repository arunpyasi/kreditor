// new dependency: ngResource is included just above
var myapp = new angular.module("myapp", ["ngResource"]);

function setLoadingScreen(show) {
    var modalLoading = UIkit.modal("#modalLoading");
    if (show) {
        modalLoading.show();
    } else {
        modalLoading.hide();
    }
}
