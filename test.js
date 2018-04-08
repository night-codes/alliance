// alert("!!!");
(function(global) {
    var alliance = global["goalliance"];
    var extender = {
        enumerable: false,
        writable: false,
        configurable: false,
        value: function(ext) {
            for (var e in ext) {
                if (ext.hasOwnProperty(e)) {
                    this[e] = ext[e];
                }
            }
        },
    }

    if (typeof alliance === "undefined") {
        alliance = {
            uris: {},
            modules: {},
            exports: {},
            circular: {},
        }
        Object.defineProperty(alliance.uris, "extend", extender);
        Object.defineProperty(alliance.modules, "extend", extender);
    }

    alliance.uris.extend({
        "main": "main.js",
        "test1": "test1.js",
        "test2": "test2/index.js",
        "test3": "test3.js",
        "test4": "test4.js",
    });
    alliance.modules.extend({
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
                }
                module.dependencies = []; 
            }  catch (err) {     
                throw '"'+alliance.uris["test4"]+'": '+err;
            }
        },
    });



    var require = function (path) {
        if (typeof path !== "string") {
            return
        }
        path = path.replace(/-/g, '_').replace(/.js$/g, '').toLowerCase();
        if (typeof alliance.exports[path] !== "undefined") {
            return alliance.exports[path];
        }

        if (alliance.circular[path]) {
            throw "Module '"+path+"' has circular dependencies!";
        }

        var module = {
            id: path,
            uri: alliance.uris[path],
            exports: {},
            dependencies: [], 
        };
        Object.defineProperty(module.exports, "extend", extender);

        if (typeof alliance.modules[path] !== "undefined") {
            alliance.circular[path] = true;
            alliance.modules[path](require, module.exports, module);
            
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
                    throw '"'+alliance.uris[path]+'": function init(): '+err
                }
            }
            alliance.exports[path] = module.exports; 
            return module.exports;
        }
    }

    if (typeof require("main") === "undefined") {
        console.error("Module 'main' is not defined")
    }

    global["require"] = require;
})(this);