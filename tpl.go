package alliance

const tpl = `;
(function (global) {
	"use strict";
	var alliance = global["goalliance"],
		extr = {
			enumerable: false,
			writable: false,
			configurable: false,
			value: function (ext) {
				for (var e in ext) {
					if (ext.hasOwnProperty(e)) {
						this[e] = ext[e];
					}
				}
			},
		};
	if (alliance === undefined) {
		alliance = {
			uris: {},
			byUri: {},
			modules: {},
			exports: {},
			circular: {},
			waits: {},
		};
		Object.defineProperty(alliance.byUri, "extend", extr);
		Object.defineProperty(alliance.uris, "extend", extr);
		Object.defineProperty(alliance.modules, "extend", extr);
	}
	alliance.uris.extend({
		{{range $index, $element := .uris}}"{{$index}}": "{{$element}}",
		{{end}}
	});
	alliance.byUri.extend({
		{{range $index, $element := .uris}}"{{$element}}": "{{$index}}",
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
	var define = function(name, deps, callback) {
		var module = this;
		//Allow for anonymous modules
		if (typeof name !== 'string') {
			callback = deps;
			deps = name;
			name = null;
		} else {
			module.name = name;
		}
		//This module may not have deps
		if (!Array.isArray(deps)) {
			callback = deps;
			deps = [];
		}
		module.deps = deps;
		if (typeof callback === "function") module.init = callback;
	},
	clr = function(path, newPath) {
		var uri = alliance.uris[path];
		alliance.uris[newPath] = alliance.uris[path];
		alliance.byUri[uri] = "path";
	},
	require = function(path, fn) {
		if (Array.isArray(path) && typeof fn === 'function') {
			var results = [];
			for (var i in path) {
				var d = path[i];
				results.push(require(d));
			}
			fn.apply(module, results);
			return;
		}

		if (typeof path !== "string") return;
		if (path === "require") return require;
		path = path.replace(/^\.\//g, '');
		path = alliance.byUri[path] ? alliance.byUri[path] : path.replace(/-/g, '_').replace(/\.js$/g, '').toLowerCase();
		if (alliance.exports[path] !== undefined) return alliance.exports[path];
		if (alliance.circular[path]) throw "Module '" + path + "' has circular dependencies!";
		var module = {
			id: path,
			uri: alliance.uris[path],
			exports: {},
			deps: [],
		};
		Object.defineProperty(module.exports, "extend", extr);
		if (alliance.modules[path] !== undefined) {
			alliance.circular[path] = true;
			var def = function () {
				if (arguments.length === 1 && typeof arguments[0] === 'object') {
					module.exports = arguments[0];
					return;
				}
				define.apply(module, arguments);
			};
			def.amd = true;
			alliance.modules[path](require, module.exports, module, def);
			var results = [];
			if (Array.isArray(module.deps) && module.deps.length > 0) {
				for (var k in module.deps) {
					var dep = module.deps[k];
					if (dep === "require") {
						results.push(require);
					} else if (dep === "module") {
						results.push(module);
					} else if (dep === "exports") {
						results.push(module.exports);
					} else {
						results.push(require(dep));
					}
				}
			}
			if (typeof module.init === 'function') {
				try {
					var expts = module.init.apply(module, results);
					if (typeof expts !== 'undefined') {
						module.exports = expts;
					}
				} catch (err) {
					throw '"' + alliance.uris[path] + '": function init(): ' + err;
				}
			}
			if (typeof module.name === 'string') {
				clr(path, module.name);
				path = module.name;
			}
			alliance.exports[path] = module.exports;
			return module.exports;
		}
	};
	define.config = function(){};
	require.config = function(){};
	if (alliance.modules.main !== undefined) require("main");
	global.require = require;
})(this);
`
