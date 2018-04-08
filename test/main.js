module.dependencies = ["dogs"];
module.init = function() {
    console.log("init main");
}
exports.test = function(){
    require("cats");
}