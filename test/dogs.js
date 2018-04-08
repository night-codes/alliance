// module.dependencies = ["sub/birds"];
module.init = function() {
    console.log("init dogs");
}
exports.print = function(){
    var cats = require("cats");
    console.log("dogs");
}
exports.gav = function(){
    console.log("gav gav!");
}