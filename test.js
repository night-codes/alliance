// alert("!!!");
(function(global) {

    var uris = {
        "main": "main.js",
        "test1": "test1.js",
        "test2": "test2/index.js",
        "test3": "test3.js",
        "test4": "test4.js",
    };
    var modules = {
        "main": function(require, exports, module, define) {
            module.dependencies = ["test1"];
            module.init = function() {
                console.log("init main");
            }
            exports.test = function(){
                require("test4");
            }
        },
        "test1": function(require, exports, module, define) {
            module.init = function() {
                require("test2");
                console.log("init test1");
            }
            module.dependencies = ["test3"];
        },
        "test2": function(require, exports, module, define) {
            module.init = function() {
                console.log("init test2");
            }
        },
        "test3": function(require, exports, module, define) {
            module.init = function() {
                console.log("init test3");
            }
            module.dependencies = ["test4"];
        },
        "test4": function(require, exports, module, define) {
            try {
                module.init = function() {
                    console.log("init test4");
                    require("test1");
                }
                module.dependencies = []; 
                require("main");
            }  catch (err) {     
                throw '"'+uris["test4"]+'": '+err.stack;
            }
        },
    }

    var exports = {};
    var circular = {};

    var require = function (path) {
        if (typeof path !== "string") {
            return
        }
        path = path.replace(/-/g, '_').replace(/.js$/g, '').toLowerCase();
        if (exports[path] !== undefined) {
            return exports[path];
        }

        if (circular[path]) {
            throw "Module '"+path+"' has circular dependencies!";
        }

        var module = {
            id: path,
            uri: uris[path],
            exports: {},
            dependencies: [], 
        };
        Object.defineProperty(module.exports, "extend", {
            enumerable: false,
            writable: false,
            configurable: false,
            value: function(exports) {
                for (var e in exports) {
                    if (exports.hasOwnProperty(e)) {
                        this[e] = exports[e];
                    }
                }
            },
        });

        if (modules[path] !== undefined) {
            circular[path] = true;
            modules[path](require, module.exports, module);
            
            if (Array.isArray(module.dependencies) && module.dependencies.length>0) {
                for (var i in module.dependencies) {
                    var d =  module.dependencies[i];
                    require(d);
                }
            }
            
            if (typeof module.init === 'function') {
                try {
                    module.init();   
                }  catch (err) {
                    throw '"'+uris[path]+'": function init(): '+err.stack
                }
            }
            exports[path] = module.exports; 
            return module.exports;
        }
    }

    if (require("main") === undefined) {
        console.error("Module 'main' is not defined")
    }

    global["require"] = require;
})(this);