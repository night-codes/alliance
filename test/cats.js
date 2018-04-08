module.dependencies = ["sub/birds"];
module.init = function() {
    console.log("init cats");
}
exports.print = function(){
    console.log("cats");
}
exports.miow = function(){
    console.log("miow miow!");
}