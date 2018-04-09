package alliance

const tpl = `
;(function(global) {
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

    if (alliance === undefined) {
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
        {{range $index, $element := .uris}}"{{$index}}": "{{$element}}",
        {{end}}
    });
    alliance.modules.extend({
        {{range $index, $element := .modules}}
        "{{$index}}": function(require, exports, module, define) {
            "use strict";
            try {
                {{$element}}
            }  catch (err) {     
                throw '"'+alliance.uris["{{$index}}"]+'": '+err;
            }  
        },     
        {{end}}
    });



    var require = function (path) {
        if (typeof path !== "string") return;
        path = path.replace(/-/g, '_').replace(/.js$/g, '').toLowerCase();
        if (alliance.exports[path] !== undefined) return alliance.exports[path];
        if (alliance.circular[path]) throw "Module '"+path+"' has circular dependencies!";

        var module = {
            id: path,
            uri: alliance.uris[path],
            exports: {},
            dependencies: [], 
        };
        Object.defineProperty(module.exports, "extend", extender);

        if (alliance.modules[path] !== undefined) {
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

    if (alliance.modules["main"] !== undefined) require("main");
    global["require"] = require;
})(this);
`
