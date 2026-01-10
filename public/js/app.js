/******/ (function(modules) { // webpackBootstrap
/******/ 	// The module cache
/******/ 	var installedModules = {};
/******/
/******/ 	// The require function
/******/ 	function __webpack_require__(moduleId) {
/******/
/******/ 		// Check if module is in cache
/******/ 		if(installedModules[moduleId]) {
/******/ 			return installedModules[moduleId].exports;
/******/ 		}
/******/ 		// Create a new module (and put it into the cache)
/******/ 		var module = installedModules[moduleId] = {
/******/ 			i: moduleId,
/******/ 			l: false,
/******/ 			exports: {}
/******/ 		};
/******/
/******/ 		// Execute the module function
/******/ 		modules[moduleId].call(module.exports, module, module.exports, __webpack_require__);
/******/
/******/ 		// Flag the module as loaded
/******/ 		module.l = true;
/******/
/******/ 		// Return the exports of the module
/******/ 		return module.exports;
/******/ 	}
/******/
/******/
/******/ 	// expose the modules object (__webpack_modules__)
/******/ 	__webpack_require__.m = modules;
/******/
/******/ 	// expose the module cache
/******/ 	__webpack_require__.c = installedModules;
/******/
/******/ 	// define getter function for harmony exports
/******/ 	__webpack_require__.d = function(exports, name, getter) {
/******/ 		if(!__webpack_require__.o(exports, name)) {
/******/ 			Object.defineProperty(exports, name, { enumerable: true, get: getter });
/******/ 		}
/******/ 	};
/******/
/******/ 	// define __esModule on exports
/******/ 	__webpack_require__.r = function(exports) {
/******/ 		if(typeof Symbol !== 'undefined' && Symbol.toStringTag) {
/******/ 			Object.defineProperty(exports, Symbol.toStringTag, { value: 'Module' });
/******/ 		}
/******/ 		Object.defineProperty(exports, '__esModule', { value: true });
/******/ 	};
/******/
/******/ 	// create a fake namespace object
/******/ 	// mode & 1: value is a module id, require it
/******/ 	// mode & 2: merge all properties of value into the ns
/******/ 	// mode & 4: return value when already ns object
/******/ 	// mode & 8|1: behave like require
/******/ 	__webpack_require__.t = function(value, mode) {
/******/ 		if(mode & 1) value = __webpack_require__(value);
/******/ 		if(mode & 8) return value;
/******/ 		if((mode & 4) && typeof value === 'object' && value && value.__esModule) return value;
/******/ 		var ns = Object.create(null);
/******/ 		__webpack_require__.r(ns);
/******/ 		Object.defineProperty(ns, 'default', { enumerable: true, value: value });
/******/ 		if(mode & 2 && typeof value != 'string') for(var key in value) __webpack_require__.d(ns, key, function(key) { return value[key]; }.bind(null, key));
/******/ 		return ns;
/******/ 	};
/******/
/******/ 	// getDefaultExport function for compatibility with non-harmony modules
/******/ 	__webpack_require__.n = function(module) {
/******/ 		var getter = module && module.__esModule ?
/******/ 			function getDefault() { return module['default']; } :
/******/ 			function getModuleExports() { return module; };
/******/ 		__webpack_require__.d(getter, 'a', getter);
/******/ 		return getter;
/******/ 	};
/******/
/******/ 	// Object.prototype.hasOwnProperty.call
/******/ 	__webpack_require__.o = function(object, property) { return Object.prototype.hasOwnProperty.call(object, property); };
/******/
/******/ 	// __webpack_public_path__
/******/ 	__webpack_require__.p = "/";
/******/
/******/
/******/ 	// Load entry module and return exports
/******/ 	return __webpack_require__(__webpack_require__.s = 0);
/******/ })
/************************************************************************/
/******/ ({

/***/ "./node_modules/babel-loader/lib/index.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/CustomerManager.vue?vue&type=script&lang=js":
/*!********************************************************************************************************************************************************************************!*\
  !*** ./node_modules/babel-loader/lib??ref--4-0!./node_modules/vue-loader/lib??vue-loader-options!./resources/js/components/salev2/CustomerManager.vue?vue&type=script&lang=js ***!
  \********************************************************************************************************************************************************************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
function _regenerator() { /*! regenerator-runtime -- Copyright (c) 2014-present, Facebook, Inc. -- license (MIT): https://github.com/babel/babel/blob/main/packages/babel-helpers/LICENSE */ var e, t, r = "function" == typeof Symbol ? Symbol : {}, n = r.iterator || "@@iterator", o = r.toStringTag || "@@toStringTag"; function i(r, n, o, i) { var c = n && n.prototype instanceof Generator ? n : Generator, u = Object.create(c.prototype); return _regeneratorDefine2(u, "_invoke", function (r, n, o) { var i, c, u, f = 0, p = o || [], y = !1, G = { p: 0, n: 0, v: e, a: d, f: d.bind(e, 4), d: function d(t, r) { return i = t, c = 0, u = e, G.n = r, a; } }; function d(r, n) { for (c = r, u = n, t = 0; !y && f && !o && t < p.length; t++) { var o, i = p[t], d = G.p, l = i[2]; r > 3 ? (o = l === n) && (u = i[(c = i[4]) ? 5 : (c = 3, 3)], i[4] = i[5] = e) : i[0] <= d && ((o = r < 2 && d < i[1]) ? (c = 0, G.v = n, G.n = i[1]) : d < l && (o = r < 3 || i[0] > n || n > l) && (i[4] = r, i[5] = n, G.n = l, c = 0)); } if (o || r > 1) return a; throw y = !0, n; } return function (o, p, l) { if (f > 1) throw TypeError("Generator is already running"); for (y && 1 === p && d(p, l), c = p, u = l; (t = c < 2 ? e : u) || !y;) { i || (c ? c < 3 ? (c > 1 && (G.n = -1), d(c, u)) : G.n = u : G.v = u); try { if (f = 2, i) { if (c || (o = "next"), t = i[o]) { if (!(t = t.call(i, u))) throw TypeError("iterator result is not an object"); if (!t.done) return t; u = t.value, c < 2 && (c = 0); } else 1 === c && (t = i["return"]) && t.call(i), c < 2 && (u = TypeError("The iterator does not provide a '" + o + "' method"), c = 1); i = e; } else if ((t = (y = G.n < 0) ? u : r.call(n, G)) !== a) break; } catch (t) { i = e, c = 1, u = t; } finally { f = 1; } } return { value: t, done: y }; }; }(r, o, i), !0), u; } var a = {}; function Generator() {} function GeneratorFunction() {} function GeneratorFunctionPrototype() {} t = Object.getPrototypeOf; var c = [][n] ? t(t([][n]())) : (_regeneratorDefine2(t = {}, n, function () { return this; }), t), u = GeneratorFunctionPrototype.prototype = Generator.prototype = Object.create(c); function f(e) { return Object.setPrototypeOf ? Object.setPrototypeOf(e, GeneratorFunctionPrototype) : (e.__proto__ = GeneratorFunctionPrototype, _regeneratorDefine2(e, o, "GeneratorFunction")), e.prototype = Object.create(u), e; } return GeneratorFunction.prototype = GeneratorFunctionPrototype, _regeneratorDefine2(u, "constructor", GeneratorFunctionPrototype), _regeneratorDefine2(GeneratorFunctionPrototype, "constructor", GeneratorFunction), GeneratorFunction.displayName = "GeneratorFunction", _regeneratorDefine2(GeneratorFunctionPrototype, o, "GeneratorFunction"), _regeneratorDefine2(u), _regeneratorDefine2(u, o, "Generator"), _regeneratorDefine2(u, n, function () { return this; }), _regeneratorDefine2(u, "toString", function () { return "[object Generator]"; }), (_regenerator = function _regenerator() { return { w: i, m: f }; })(); }
function _regeneratorDefine2(e, r, n, t) { var i = Object.defineProperty; try { i({}, "", {}); } catch (e) { i = 0; } _regeneratorDefine2 = function _regeneratorDefine(e, r, n, t) { function o(r, n) { _regeneratorDefine2(e, r, function (e) { return this._invoke(r, n, e); }); } r ? i ? i(e, r, { value: n, enumerable: !t, configurable: !t, writable: !t }) : e[r] = n : (o("next", 0), o("throw", 1), o("return", 2)); }, _regeneratorDefine2(e, r, n, t); }
function asyncGeneratorStep(n, t, e, r, o, a, c) { try { var i = n[a](c), u = i.value; } catch (n) { return void e(n); } i.done ? t(u) : Promise.resolve(u).then(r, o); }
function _asyncToGenerator(n) { return function () { var t = this, e = arguments; return new Promise(function (r, o) { var a = n.apply(t, e); function _next(n) { asyncGeneratorStep(a, r, o, _next, _throw, "next", n); } function _throw(n) { asyncGeneratorStep(a, r, o, _next, _throw, "throw", n); } _next(void 0); }); }; }
/* harmony default export */ __webpack_exports__["default"] = ({
  name: 'CustomerManager',
  data: function data() {
    return {
      consumers: [],
      loading: false,
      search: '',
      page: 1,
      total: 0
    };
  },
  mounted: function mounted() {
    this.fetchConsumers();
  },
  methods: {
    fetchConsumers: function fetchConsumers() {
      var _arguments = arguments,
        _this = this;
      return _asyncToGenerator(/*#__PURE__*/_regenerator().m(function _callee() {
        var page, res, _t;
        return _regenerator().w(function (_context) {
          while (1) switch (_context.p = _context.n) {
            case 0:
              page = _arguments.length > 0 && _arguments[0] !== undefined ? _arguments[0] : 1;
              _this.page = page;
              _this.loading = true;
              _context.p = 1;
              _context.n = 2;
              return axios.get('/salev2/api/consumers', {
                params: {
                  page: page,
                  search: _this.search
                }
              });
            case 2:
              res = _context.v;
              _this.consumers = res.data.data;
              _this.total = res.data.total;
              _context.n = 4;
              break;
            case 3:
              _context.p = 3;
              _t = _context.v;
              console.error(_t);
              _this.$message.error('Failed to load consumers');
            case 4:
              _context.p = 4;
              _this.loading = false;
              return _context.f(4);
            case 5:
              return _context.a(2);
          }
        }, _callee, null, [[1, 3, 4, 5]]);
      }))();
    },
    debouncedFetch: _.debounce(function () {
      this.fetchConsumers(1);
    }, 500),
    createOrder: function createOrder(customer) {
      this.$router.push({
        name: 'order-create',
        query: {
          consumer_id: customer.id
        }
      });
    }
  }
});

/***/ }),

/***/ "./node_modules/babel-loader/lib/index.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/Dashboard.vue?vue&type=script&lang=js":
/*!**************************************************************************************************************************************************************************!*\
  !*** ./node_modules/babel-loader/lib??ref--4-0!./node_modules/vue-loader/lib??vue-loader-options!./resources/js/components/salev2/Dashboard.vue?vue&type=script&lang=js ***!
  \**************************************************************************************************************************************************************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
function _toConsumableArray(r) { return _arrayWithoutHoles(r) || _iterableToArray(r) || _unsupportedIterableToArray(r) || _nonIterableSpread(); }
function _nonIterableSpread() { throw new TypeError("Invalid attempt to spread non-iterable instance.\nIn order to be iterable, non-array objects must have a [Symbol.iterator]() method."); }
function _unsupportedIterableToArray(r, a) { if (r) { if ("string" == typeof r) return _arrayLikeToArray(r, a); var t = {}.toString.call(r).slice(8, -1); return "Object" === t && r.constructor && (t = r.constructor.name), "Map" === t || "Set" === t ? Array.from(r) : "Arguments" === t || /^(?:Ui|I)nt(?:8|16|32)(?:Clamped)?Array$/.test(t) ? _arrayLikeToArray(r, a) : void 0; } }
function _iterableToArray(r) { if ("undefined" != typeof Symbol && null != r[Symbol.iterator] || null != r["@@iterator"]) return Array.from(r); }
function _arrayWithoutHoles(r) { if (Array.isArray(r)) return _arrayLikeToArray(r); }
function _arrayLikeToArray(r, a) { (null == a || a > r.length) && (a = r.length); for (var e = 0, n = Array(a); e < a; e++) n[e] = r[e]; return n; }
function _regenerator() { /*! regenerator-runtime -- Copyright (c) 2014-present, Facebook, Inc. -- license (MIT): https://github.com/babel/babel/blob/main/packages/babel-helpers/LICENSE */ var e, t, r = "function" == typeof Symbol ? Symbol : {}, n = r.iterator || "@@iterator", o = r.toStringTag || "@@toStringTag"; function i(r, n, o, i) { var c = n && n.prototype instanceof Generator ? n : Generator, u = Object.create(c.prototype); return _regeneratorDefine2(u, "_invoke", function (r, n, o) { var i, c, u, f = 0, p = o || [], y = !1, G = { p: 0, n: 0, v: e, a: d, f: d.bind(e, 4), d: function d(t, r) { return i = t, c = 0, u = e, G.n = r, a; } }; function d(r, n) { for (c = r, u = n, t = 0; !y && f && !o && t < p.length; t++) { var o, i = p[t], d = G.p, l = i[2]; r > 3 ? (o = l === n) && (u = i[(c = i[4]) ? 5 : (c = 3, 3)], i[4] = i[5] = e) : i[0] <= d && ((o = r < 2 && d < i[1]) ? (c = 0, G.v = n, G.n = i[1]) : d < l && (o = r < 3 || i[0] > n || n > l) && (i[4] = r, i[5] = n, G.n = l, c = 0)); } if (o || r > 1) return a; throw y = !0, n; } return function (o, p, l) { if (f > 1) throw TypeError("Generator is already running"); for (y && 1 === p && d(p, l), c = p, u = l; (t = c < 2 ? e : u) || !y;) { i || (c ? c < 3 ? (c > 1 && (G.n = -1), d(c, u)) : G.n = u : G.v = u); try { if (f = 2, i) { if (c || (o = "next"), t = i[o]) { if (!(t = t.call(i, u))) throw TypeError("iterator result is not an object"); if (!t.done) return t; u = t.value, c < 2 && (c = 0); } else 1 === c && (t = i["return"]) && t.call(i), c < 2 && (u = TypeError("The iterator does not provide a '" + o + "' method"), c = 1); i = e; } else if ((t = (y = G.n < 0) ? u : r.call(n, G)) !== a) break; } catch (t) { i = e, c = 1, u = t; } finally { f = 1; } } return { value: t, done: y }; }; }(r, o, i), !0), u; } var a = {}; function Generator() {} function GeneratorFunction() {} function GeneratorFunctionPrototype() {} t = Object.getPrototypeOf; var c = [][n] ? t(t([][n]())) : (_regeneratorDefine2(t = {}, n, function () { return this; }), t), u = GeneratorFunctionPrototype.prototype = Generator.prototype = Object.create(c); function f(e) { return Object.setPrototypeOf ? Object.setPrototypeOf(e, GeneratorFunctionPrototype) : (e.__proto__ = GeneratorFunctionPrototype, _regeneratorDefine2(e, o, "GeneratorFunction")), e.prototype = Object.create(u), e; } return GeneratorFunction.prototype = GeneratorFunctionPrototype, _regeneratorDefine2(u, "constructor", GeneratorFunctionPrototype), _regeneratorDefine2(GeneratorFunctionPrototype, "constructor", GeneratorFunction), GeneratorFunction.displayName = "GeneratorFunction", _regeneratorDefine2(GeneratorFunctionPrototype, o, "GeneratorFunction"), _regeneratorDefine2(u), _regeneratorDefine2(u, o, "Generator"), _regeneratorDefine2(u, n, function () { return this; }), _regeneratorDefine2(u, "toString", function () { return "[object Generator]"; }), (_regenerator = function _regenerator() { return { w: i, m: f }; })(); }
function _regeneratorDefine2(e, r, n, t) { var i = Object.defineProperty; try { i({}, "", {}); } catch (e) { i = 0; } _regeneratorDefine2 = function _regeneratorDefine(e, r, n, t) { function o(r, n) { _regeneratorDefine2(e, r, function (e) { return this._invoke(r, n, e); }); } r ? i ? i(e, r, { value: n, enumerable: !t, configurable: !t, writable: !t }) : e[r] = n : (o("next", 0), o("throw", 1), o("return", 2)); }, _regeneratorDefine2(e, r, n, t); }
function asyncGeneratorStep(n, t, e, r, o, a, c) { try { var i = n[a](c), u = i.value; } catch (n) { return void e(n); } i.done ? t(u) : Promise.resolve(u).then(r, o); }
function _asyncToGenerator(n) { return function () { var t = this, e = arguments; return new Promise(function (r, o) { var a = n.apply(t, e); function _next(n) { asyncGeneratorStep(a, r, o, _next, _throw, "next", n); } function _throw(n) { asyncGeneratorStep(a, r, o, _next, _throw, "throw", n); } _next(void 0); }); }; }
// Using global axios
/* harmony default export */ __webpack_exports__["default"] = ({
  name: 'Dashboard',
  data: function data() {
    return {
      data: null,
      loading: false,
      productTypes: ['TROUSERS', 'JACKET', 'SHIRT', 'BLAZER', 'TIE', 'GILE', 'BELT', 'POLO SHIRT', 'HAT', 'SUIT'],
      statCards: [{
        title: 'New Orders Rev',
        value: '0đ',
        icon: 'el-icon-message-solid',
        bgClass: 'bg-orange-500'
      }, {
        title: 'New Orders Count',
        value: '0',
        icon: 'el-icon-bell',
        bgClass: 'bg-orange-400'
      }, {
        title: 'Monthly Revenue',
        value: '0đ',
        icon: 'el-icon-money',
        bgClass: 'bg-emerald-500'
      }, {
        title: 'Monthly Count',
        value: '0',
        icon: 'el-icon-shopping-cart-full',
        bgClass: 'bg-blue-500'
      }]
    };
  },
  computed: {
    allTypeStats: function allTypeStats() {
      return this.productTypes.map(function (t) {
        return {
          type: t
        };
      });
    }
  },
  mounted: function mounted() {
    this.fetchStats();
  },
  methods: {
    formatPrice: function formatPrice(val) {
      if (!val) return '0đ';
      return new Intl.NumberFormat('vi-VN', {
        style: 'currency',
        currency: 'VND',
        maximumFractionDigits: 0
      }).format(val);
    },
    fetchStats: function fetchStats() {
      var _this = this;
      return _asyncToGenerator(/*#__PURE__*/_regenerator().m(function _callee() {
        var response, _t;
        return _regenerator().w(function (_context) {
          while (1) switch (_context.p = _context.n) {
            case 0:
              _this.loading = true;
              _context.p = 1;
              _context.n = 2;
              return axios.get('/salev2/api/dashboard-stats');
            case 2:
              response = _context.v;
              _this.data = response.data;

              // Update Stat Cards (Legacy Sync)
              _this.statCards[0].value = _this.formatPrice(_this.data.new_orders_revenue);
              _this.statCards[1].value = _this.data.new_orders_count;
              _this.statCards[2].value = _this.formatPrice(_this.data.month_orders_revenue);
              _this.statCards[3].value = _this.data.month_orders_count;
              _context.n = 4;
              break;
            case 3:
              _context.p = 3;
              _t = _context.v;
              console.error('Failed to fetch stats:', _t);
            case 4:
              _context.p = 4;
              _this.loading = false;
              return _context.f(4);
            case 5:
              return _context.a(2);
          }
        }, _callee, null, [[1, 3, 4, 5]]);
      }))();
    },
    getCount: function getCount(stats, type) {
      var status = arguments.length > 2 && arguments[2] !== undefined ? arguments[2] : 'SOLD';
      if (!stats) return 0;
      var found = stats.find(function (s) {
        return s.type === type && (!s.status || s.status === status);
      });
      return found ? found.count : 0;
    },
    getTrendWidth: function getTrendWidth(amount) {
      if (!this.data || !this.data.monthly_sales.length) return 0;
      var max = Math.max.apply(Math, _toConsumableArray(this.data.monthly_sales.map(function (m) {
        return m.final_amount;
      })));
      return amount / max * 100;
    }
  }
});

/***/ }),

/***/ "./node_modules/babel-loader/lib/index.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/OrderCreate.vue?vue&type=script&lang=js":
/*!****************************************************************************************************************************************************************************!*\
  !*** ./node_modules/babel-loader/lib??ref--4-0!./node_modules/vue-loader/lib??vue-loader-options!./resources/js/components/salev2/OrderCreate.vue?vue&type=script&lang=js ***!
  \****************************************************************************************************************************************************************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
function _typeof(o) { "@babel/helpers - typeof"; return _typeof = "function" == typeof Symbol && "symbol" == typeof Symbol.iterator ? function (o) { return typeof o; } : function (o) { return o && "function" == typeof Symbol && o.constructor === Symbol && o !== Symbol.prototype ? "symbol" : typeof o; }, _typeof(o); }
function ownKeys(e, r) { var t = Object.keys(e); if (Object.getOwnPropertySymbols) { var o = Object.getOwnPropertySymbols(e); r && (o = o.filter(function (r) { return Object.getOwnPropertyDescriptor(e, r).enumerable; })), t.push.apply(t, o); } return t; }
function _objectSpread(e) { for (var r = 1; r < arguments.length; r++) { var t = null != arguments[r] ? arguments[r] : {}; r % 2 ? ownKeys(Object(t), !0).forEach(function (r) { _defineProperty(e, r, t[r]); }) : Object.getOwnPropertyDescriptors ? Object.defineProperties(e, Object.getOwnPropertyDescriptors(t)) : ownKeys(Object(t)).forEach(function (r) { Object.defineProperty(e, r, Object.getOwnPropertyDescriptor(t, r)); }); } return e; }
function _defineProperty(e, r, t) { return (r = _toPropertyKey(r)) in e ? Object.defineProperty(e, r, { value: t, enumerable: !0, configurable: !0, writable: !0 }) : e[r] = t, e; }
function _toPropertyKey(t) { var i = _toPrimitive(t, "string"); return "symbol" == _typeof(i) ? i : i + ""; }
function _toPrimitive(t, r) { if ("object" != _typeof(t) || !t) return t; var e = t[Symbol.toPrimitive]; if (void 0 !== e) { var i = e.call(t, r || "default"); if ("object" != _typeof(i)) return i; throw new TypeError("@@toPrimitive must return a primitive value."); } return ("string" === r ? String : Number)(t); }
function _regenerator() { /*! regenerator-runtime -- Copyright (c) 2014-present, Facebook, Inc. -- license (MIT): https://github.com/babel/babel/blob/main/packages/babel-helpers/LICENSE */ var e, t, r = "function" == typeof Symbol ? Symbol : {}, n = r.iterator || "@@iterator", o = r.toStringTag || "@@toStringTag"; function i(r, n, o, i) { var c = n && n.prototype instanceof Generator ? n : Generator, u = Object.create(c.prototype); return _regeneratorDefine2(u, "_invoke", function (r, n, o) { var i, c, u, f = 0, p = o || [], y = !1, G = { p: 0, n: 0, v: e, a: d, f: d.bind(e, 4), d: function d(t, r) { return i = t, c = 0, u = e, G.n = r, a; } }; function d(r, n) { for (c = r, u = n, t = 0; !y && f && !o && t < p.length; t++) { var o, i = p[t], d = G.p, l = i[2]; r > 3 ? (o = l === n) && (u = i[(c = i[4]) ? 5 : (c = 3, 3)], i[4] = i[5] = e) : i[0] <= d && ((o = r < 2 && d < i[1]) ? (c = 0, G.v = n, G.n = i[1]) : d < l && (o = r < 3 || i[0] > n || n > l) && (i[4] = r, i[5] = n, G.n = l, c = 0)); } if (o || r > 1) return a; throw y = !0, n; } return function (o, p, l) { if (f > 1) throw TypeError("Generator is already running"); for (y && 1 === p && d(p, l), c = p, u = l; (t = c < 2 ? e : u) || !y;) { i || (c ? c < 3 ? (c > 1 && (G.n = -1), d(c, u)) : G.n = u : G.v = u); try { if (f = 2, i) { if (c || (o = "next"), t = i[o]) { if (!(t = t.call(i, u))) throw TypeError("iterator result is not an object"); if (!t.done) return t; u = t.value, c < 2 && (c = 0); } else 1 === c && (t = i["return"]) && t.call(i), c < 2 && (u = TypeError("The iterator does not provide a '" + o + "' method"), c = 1); i = e; } else if ((t = (y = G.n < 0) ? u : r.call(n, G)) !== a) break; } catch (t) { i = e, c = 1, u = t; } finally { f = 1; } } return { value: t, done: y }; }; }(r, o, i), !0), u; } var a = {}; function Generator() {} function GeneratorFunction() {} function GeneratorFunctionPrototype() {} t = Object.getPrototypeOf; var c = [][n] ? t(t([][n]())) : (_regeneratorDefine2(t = {}, n, function () { return this; }), t), u = GeneratorFunctionPrototype.prototype = Generator.prototype = Object.create(c); function f(e) { return Object.setPrototypeOf ? Object.setPrototypeOf(e, GeneratorFunctionPrototype) : (e.__proto__ = GeneratorFunctionPrototype, _regeneratorDefine2(e, o, "GeneratorFunction")), e.prototype = Object.create(u), e; } return GeneratorFunction.prototype = GeneratorFunctionPrototype, _regeneratorDefine2(u, "constructor", GeneratorFunctionPrototype), _regeneratorDefine2(GeneratorFunctionPrototype, "constructor", GeneratorFunction), GeneratorFunction.displayName = "GeneratorFunction", _regeneratorDefine2(GeneratorFunctionPrototype, o, "GeneratorFunction"), _regeneratorDefine2(u), _regeneratorDefine2(u, o, "Generator"), _regeneratorDefine2(u, n, function () { return this; }), _regeneratorDefine2(u, "toString", function () { return "[object Generator]"; }), (_regenerator = function _regenerator() { return { w: i, m: f }; })(); }
function _regeneratorDefine2(e, r, n, t) { var i = Object.defineProperty; try { i({}, "", {}); } catch (e) { i = 0; } _regeneratorDefine2 = function _regeneratorDefine(e, r, n, t) { function o(r, n) { _regeneratorDefine2(e, r, function (e) { return this._invoke(r, n, e); }); } r ? i ? i(e, r, { value: n, enumerable: !t, configurable: !t, writable: !t }) : e[r] = n : (o("next", 0), o("throw", 1), o("return", 2)); }, _regeneratorDefine2(e, r, n, t); }
function asyncGeneratorStep(n, t, e, r, o, a, c) { try { var i = n[a](c), u = i.value; } catch (n) { return void e(n); } i.done ? t(u) : Promise.resolve(u).then(r, o); }
function _asyncToGenerator(n) { return function () { var t = this, e = arguments; return new Promise(function (r, o) { var a = n.apply(t, e); function _next(n) { asyncGeneratorStep(a, r, o, _next, _throw, "next", n); } function _throw(n) { asyncGeneratorStep(a, r, o, _next, _throw, "throw", n); } _next(void 0); }); }; }
/* harmony default export */ __webpack_exports__["default"] = ({
  name: 'OrderCreate',
  data: function data() {
    return {
      customer: {
        name: '',
        phone: '',
        address: ''
      },
      order: {
        note: '',
        desc: '',
        additional_amount: 0,
        freeship_amount: 0,
        deposit_amount: 0,
        is_paid: false
      },
      productSearch: '',
      searchResults: [],
      selectedProducts: [],
      loadingCustomer: false,
      loadingProducts: false,
      submitting: false
    };
  },
  computed: {
    subtotal: function subtotal() {
      return this.selectedProducts.reduce(function (sum, p) {
        return sum + Number(p.price);
      }, 0);
    },
    finalTotal: function finalTotal() {
      var total = this.subtotal + this.order.additional_amount - this.order.freeship_amount;
      return total > 0 ? total : 0;
    }
  },
  watch: {
    '$route.query.consumer_id': {
      immediate: true,
      handler: function handler(id) {
        if (id) this.fetchConsumer(id);
      }
    }
  },
  methods: {
    formatPrice: function formatPrice(val) {
      return new Intl.NumberFormat('vi-VN', {
        style: 'currency',
        currency: 'VND',
        maximumFractionDigits: 0
      }).format(val);
    },
    fetchConsumer: function fetchConsumer(id) {
      var _this = this;
      return _asyncToGenerator(/*#__PURE__*/_regenerator().m(function _callee() {
        var res, _t;
        return _regenerator().w(function (_context) {
          while (1) switch (_context.p = _context.n) {
            case 0:
              _this.loadingCustomer = true;
              _context.p = 1;
              _context.n = 2;
              return axios.get("/salev2/api/consumer/".concat(id));
            case 2:
              res = _context.v;
              _this.customer = {
                name: res.data.name,
                phone: res.data.phone,
                address: res.data.address
              };
              _context.n = 4;
              break;
            case 3:
              _context.p = 3;
              _t = _context.v;
              _this.$message.error('Failed to load customer');
            case 4:
              _context.p = 4;
              _this.loadingCustomer = false;
              return _context.f(4);
            case 5:
              return _context.a(2);
          }
        }, _callee, null, [[1, 3, 4, 5]]);
      }))();
    },
    debouncedSearch: _.debounce(function () {
      this.searchProducts();
    }, 500),
    searchProducts: function searchProducts() {
      var _this2 = this;
      return _asyncToGenerator(/*#__PURE__*/_regenerator().m(function _callee2() {
        var res, _t2;
        return _regenerator().w(function (_context2) {
          while (1) switch (_context2.p = _context2.n) {
            case 0:
              if (_this2.productSearch) {
                _context2.n = 1;
                break;
              }
              return _context2.a(2);
            case 1:
              _this2.loadingProducts = true;
              _context2.p = 2;
              _context2.n = 3;
              return axios.get('/salev2/api/products', {
                params: {
                  search: _this2.productSearch,
                  status: 'AVAILABLE'
                }
              });
            case 3:
              res = _context2.v;
              _this2.searchResults = res.data.data;
              _context2.n = 5;
              break;
            case 4:
              _context2.p = 4;
              _t2 = _context2.v;
              console.error(_t2);
            case 5:
              _context2.p = 5;
              _this2.loadingProducts = false;
              return _context2.f(5);
            case 6:
              return _context2.a(2);
          }
        }, _callee2, null, [[2, 4, 5, 6]]);
      }))();
    },
    toggleProduct: function toggleProduct(prod) {
      var idx = this.selectedProducts.findIndex(function (p) {
        return p.id === prod.id;
      });
      if (idx > -1) {
        this.selectedProducts.splice(idx, 1);
      } else {
        this.selectedProducts.push(prod);
      }
    },
    isSelected: function isSelected(prod) {
      return this.selectedProducts.some(function (p) {
        return p.id === prod.id;
      });
    },
    submitOrder: function submitOrder() {
      var _this3 = this;
      return _asyncToGenerator(/*#__PURE__*/_regenerator().m(function _callee3() {
        var payload, _t3;
        return _regenerator().w(function (_context3) {
          while (1) switch (_context3.p = _context3.n) {
            case 0:
              if (!(!_this3.customer.name || !_this3.customer.phone)) {
                _context3.n = 1;
                break;
              }
              return _context3.a(2, _this3.$message.error('Customer info required'));
            case 1:
              if (_this3.selectedProducts.length) {
                _context3.n = 2;
                break;
              }
              return _context3.a(2, _this3.$message.error('No products selected'));
            case 2:
              _this3.submitting = true;
              _context3.p = 3;
              payload = _objectSpread(_objectSpread(_objectSpread({}, _this3.customer), {}, {
                // Flatten customer fields
                consumer_id: _this3.$route.query.consumer_id
              }, _this3.order), {}, {
                total_amount: _this3.finalTotal,
                // Legacy expects total_amount as Final? No, let's allow backend to calculate or send calc.
                // Actually storeOrder expects: products (array of IDs)
                products: _this3.selectedProducts.map(function (p) {
                  return p.id;
                }),
                products_info: _this3.selectedProducts.map(function (p) {
                  return {
                    id: p.id,
                    price: p.price
                  };
                })
              });
              _context3.n = 4;
              return axios.post('/salev2/api/order', payload);
            case 4:
              _this3.$message.success('Order Created!');
              _this3.$router.push('/salev2/orders');
              _context3.n = 6;
              break;
            case 5:
              _context3.p = 5;
              _t3 = _context3.v;
              _this3.$message.error('Failed to create order');
            case 6:
              _context3.p = 6;
              _this3.submitting = false;
              return _context3.f(6);
            case 7:
              return _context3.a(2);
          }
        }, _callee3, null, [[3, 5, 6, 7]]);
      }))();
    }
  }
});

/***/ }),

/***/ "./node_modules/babel-loader/lib/index.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/OrderEdit.vue?vue&type=script&lang=js":
/*!**************************************************************************************************************************************************************************!*\
  !*** ./node_modules/babel-loader/lib??ref--4-0!./node_modules/vue-loader/lib??vue-loader-options!./resources/js/components/salev2/OrderEdit.vue?vue&type=script&lang=js ***!
  \**************************************************************************************************************************************************************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
function _typeof(o) { "@babel/helpers - typeof"; return _typeof = "function" == typeof Symbol && "symbol" == typeof Symbol.iterator ? function (o) { return typeof o; } : function (o) { return o && "function" == typeof Symbol && o.constructor === Symbol && o !== Symbol.prototype ? "symbol" : typeof o; }, _typeof(o); }
function ownKeys(e, r) { var t = Object.keys(e); if (Object.getOwnPropertySymbols) { var o = Object.getOwnPropertySymbols(e); r && (o = o.filter(function (r) { return Object.getOwnPropertyDescriptor(e, r).enumerable; })), t.push.apply(t, o); } return t; }
function _objectSpread(e) { for (var r = 1; r < arguments.length; r++) { var t = null != arguments[r] ? arguments[r] : {}; r % 2 ? ownKeys(Object(t), !0).forEach(function (r) { _defineProperty(e, r, t[r]); }) : Object.getOwnPropertyDescriptors ? Object.defineProperties(e, Object.getOwnPropertyDescriptors(t)) : ownKeys(Object(t)).forEach(function (r) { Object.defineProperty(e, r, Object.getOwnPropertyDescriptor(t, r)); }); } return e; }
function _defineProperty(e, r, t) { return (r = _toPropertyKey(r)) in e ? Object.defineProperty(e, r, { value: t, enumerable: !0, configurable: !0, writable: !0 }) : e[r] = t, e; }
function _toPropertyKey(t) { var i = _toPrimitive(t, "string"); return "symbol" == _typeof(i) ? i : i + ""; }
function _toPrimitive(t, r) { if ("object" != _typeof(t) || !t) return t; var e = t[Symbol.toPrimitive]; if (void 0 !== e) { var i = e.call(t, r || "default"); if ("object" != _typeof(i)) return i; throw new TypeError("@@toPrimitive must return a primitive value."); } return ("string" === r ? String : Number)(t); }
function _regenerator() { /*! regenerator-runtime -- Copyright (c) 2014-present, Facebook, Inc. -- license (MIT): https://github.com/babel/babel/blob/main/packages/babel-helpers/LICENSE */ var e, t, r = "function" == typeof Symbol ? Symbol : {}, n = r.iterator || "@@iterator", o = r.toStringTag || "@@toStringTag"; function i(r, n, o, i) { var c = n && n.prototype instanceof Generator ? n : Generator, u = Object.create(c.prototype); return _regeneratorDefine2(u, "_invoke", function (r, n, o) { var i, c, u, f = 0, p = o || [], y = !1, G = { p: 0, n: 0, v: e, a: d, f: d.bind(e, 4), d: function d(t, r) { return i = t, c = 0, u = e, G.n = r, a; } }; function d(r, n) { for (c = r, u = n, t = 0; !y && f && !o && t < p.length; t++) { var o, i = p[t], d = G.p, l = i[2]; r > 3 ? (o = l === n) && (u = i[(c = i[4]) ? 5 : (c = 3, 3)], i[4] = i[5] = e) : i[0] <= d && ((o = r < 2 && d < i[1]) ? (c = 0, G.v = n, G.n = i[1]) : d < l && (o = r < 3 || i[0] > n || n > l) && (i[4] = r, i[5] = n, G.n = l, c = 0)); } if (o || r > 1) return a; throw y = !0, n; } return function (o, p, l) { if (f > 1) throw TypeError("Generator is already running"); for (y && 1 === p && d(p, l), c = p, u = l; (t = c < 2 ? e : u) || !y;) { i || (c ? c < 3 ? (c > 1 && (G.n = -1), d(c, u)) : G.n = u : G.v = u); try { if (f = 2, i) { if (c || (o = "next"), t = i[o]) { if (!(t = t.call(i, u))) throw TypeError("iterator result is not an object"); if (!t.done) return t; u = t.value, c < 2 && (c = 0); } else 1 === c && (t = i["return"]) && t.call(i), c < 2 && (u = TypeError("The iterator does not provide a '" + o + "' method"), c = 1); i = e; } else if ((t = (y = G.n < 0) ? u : r.call(n, G)) !== a) break; } catch (t) { i = e, c = 1, u = t; } finally { f = 1; } } return { value: t, done: y }; }; }(r, o, i), !0), u; } var a = {}; function Generator() {} function GeneratorFunction() {} function GeneratorFunctionPrototype() {} t = Object.getPrototypeOf; var c = [][n] ? t(t([][n]())) : (_regeneratorDefine2(t = {}, n, function () { return this; }), t), u = GeneratorFunctionPrototype.prototype = Generator.prototype = Object.create(c); function f(e) { return Object.setPrototypeOf ? Object.setPrototypeOf(e, GeneratorFunctionPrototype) : (e.__proto__ = GeneratorFunctionPrototype, _regeneratorDefine2(e, o, "GeneratorFunction")), e.prototype = Object.create(u), e; } return GeneratorFunction.prototype = GeneratorFunctionPrototype, _regeneratorDefine2(u, "constructor", GeneratorFunctionPrototype), _regeneratorDefine2(GeneratorFunctionPrototype, "constructor", GeneratorFunction), GeneratorFunction.displayName = "GeneratorFunction", _regeneratorDefine2(GeneratorFunctionPrototype, o, "GeneratorFunction"), _regeneratorDefine2(u), _regeneratorDefine2(u, o, "Generator"), _regeneratorDefine2(u, n, function () { return this; }), _regeneratorDefine2(u, "toString", function () { return "[object Generator]"; }), (_regenerator = function _regenerator() { return { w: i, m: f }; })(); }
function _regeneratorDefine2(e, r, n, t) { var i = Object.defineProperty; try { i({}, "", {}); } catch (e) { i = 0; } _regeneratorDefine2 = function _regeneratorDefine(e, r, n, t) { function o(r, n) { _regeneratorDefine2(e, r, function (e) { return this._invoke(r, n, e); }); } r ? i ? i(e, r, { value: n, enumerable: !t, configurable: !t, writable: !t }) : e[r] = n : (o("next", 0), o("throw", 1), o("return", 2)); }, _regeneratorDefine2(e, r, n, t); }
function asyncGeneratorStep(n, t, e, r, o, a, c) { try { var i = n[a](c), u = i.value; } catch (n) { return void e(n); } i.done ? t(u) : Promise.resolve(u).then(r, o); }
function _asyncToGenerator(n) { return function () { var t = this, e = arguments; return new Promise(function (r, o) { var a = n.apply(t, e); function _next(n) { asyncGeneratorStep(a, r, o, _next, _throw, "next", n); } function _throw(n) { asyncGeneratorStep(a, r, o, _next, _throw, "throw", n); } _next(void 0); }); }; }
/* harmony default export */ __webpack_exports__["default"] = ({
  name: 'OrderEdit',
  data: function data() {
    return {
      loading: false,
      submitting: false,
      order: null,
      showProductSelector: false,
      productSearch: '',
      availableProducts: [],
      additionalProds: []
    };
  },
  computed: {
    subtotal: function subtotal() {
      if (!this.order) return 0;
      return this.order.products.reduce(function (sum, p) {
        return sum + p.price;
      }, 0);
    },
    totalAmount: function totalAmount() {
      if (!this.order) return 0;
      var total = this.subtotal + (this.order.additional_amount || 0);
      if (!this.order.is_freeship) total += 20000; // Assuming 20k ship fee if not free
      // Actually legacy logic subtracts freeship_amount but here we simplify: 
      // Legacy: total_amount = products + additional - freeship.
      // But we use is_freeship flag. Let's stick to simple total calculation for display.
      // Wait, legacy stores final_amount.
      return total;
    },
    viettelAmount: function viettelAmount() {
      if (!this.order) return 0;
      if (this.order.is_paid) return 0;
      return this.totalAmount - (this.order.deposit_amount || 0);
    }
  },
  mounted: function mounted() {
    this.fetchOrder();
    this.fetchAvailableProducts();
  },
  methods: {
    formatPrice: function formatPrice(val) {
      return new Intl.NumberFormat('vi-VN', {
        style: 'currency',
        currency: 'VND',
        maximumFractionDigits: 0
      }).format(val);
    },
    fetchOrder: function fetchOrder() {
      var _this = this;
      return _asyncToGenerator(/*#__PURE__*/_regenerator().m(function _callee() {
        var id, res, _t;
        return _regenerator().w(function (_context) {
          while (1) switch (_context.p = _context.n) {
            case 0:
              _this.loading = true;
              _context.p = 1;
              id = _this.$route.params.id;
              _context.n = 2;
              return axios.get("/salev2/api/order/".concat(id));
            case 2:
              res = _context.v;
              _this.order = res.data;
              // Ensure array fields
              if (!_this.order.products) _this.order.products = [];
              _context.n = 4;
              break;
            case 3:
              _context.p = 3;
              _t = _context.v;
              _this.$message.error('Failed to load order');
              _this.$router.push('/salev2/orders');
            case 4:
              _context.p = 4;
              _this.loading = false;
              return _context.f(4);
            case 5:
              return _context.a(2);
          }
        }, _callee, null, [[1, 3, 4, 5]]);
      }))();
    },
    fetchAvailableProducts: function fetchAvailableProducts() {
      var _this2 = this;
      return _asyncToGenerator(/*#__PURE__*/_regenerator().m(function _callee2() {
        var res, _t2;
        return _regenerator().w(function (_context2) {
          while (1) switch (_context2.p = _context2.n) {
            case 0:
              _context2.p = 0;
              _context2.n = 1;
              return axios.get('/salev2/api/products', {
                params: {
                  status: 'AVAILABLE',
                  search: _this2.productSearch,
                  limit: 50
                }
              });
            case 1:
              res = _context2.v;
              _this2.availableProducts = res.data.data;
              _context2.n = 3;
              break;
            case 2:
              _context2.p = 2;
              _t2 = _context2.v;
            case 3:
              return _context2.a(2);
          }
        }, _callee2, null, [[0, 2]]);
      }))();
    },
    debouncedSearch: _.debounce(function () {
      this.fetchAvailableProducts();
    }, 500),
    isProductSelected: function isProductSelected(id) {
      return this.order && this.order.products.some(function (p) {
        return p.id === id;
      });
    },
    toggleProductSelect: function toggleProductSelect(product) {
      if (this.isProductSelected(product.id)) {
        // If it was originally in order, we remove it. 
        // Wait, if we remove it, it goes back to AVAILABLE.
        // If we add it, it goes to SOLD.
        var idx = this.order.products.findIndex(function (p) {
          return p.id === product.id;
        });
        if (idx > -1) this.order.products.splice(idx, 1);
      } else {
        this.order.products.push(product);
      }
    },
    removeProduct: function removeProduct(index) {
      this.order.products.splice(index, 1);
    },
    saveOrder: function saveOrder() {
      var _this3 = this;
      return _asyncToGenerator(/*#__PURE__*/_regenerator().m(function _callee3() {
        var payload, _t3;
        return _regenerator().w(function (_context3) {
          while (1) switch (_context3.p = _context3.n) {
            case 0:
              _this3.submitting = true;
              _context3.p = 1;
              // Prepare payload
              payload = _objectSpread(_objectSpread({}, _this3.order), {}, {
                products: _this3.order.products.map(function (p) {
                  return p.id;
                }),
                // Send price updates
                products_info: _this3.order.products.map(function (p) {
                  return {
                    id: p.id,
                    price: p.price
                  };
                }),
                // Recalculate financial fields to ensure backend gets correct numbers
                total_amount: _this3.totalAmount,
                // Backend might recalc but good to send
                final_amount: _this3.totalAmount
              });
              _context3.n = 2;
              return axios.post("/salev2/api/order/".concat(_this3.order.id, "/update"), payload);
            case 2:
              _this3.$message.success('Order Updated');
              _this3.$router.push('/salev2/orders');
              _context3.n = 4;
              break;
            case 3:
              _context3.p = 3;
              _t3 = _context3.v;
              _this3.$message.error('Update failed');
            case 4:
              _context3.p = 4;
              _this3.submitting = false;
              return _context3.f(4);
            case 5:
              return _context3.a(2);
          }
        }, _callee3, null, [[1, 3, 4, 5]]);
      }))();
    }
  }
});

/***/ }),

/***/ "./node_modules/babel-loader/lib/index.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/OrderManager.vue?vue&type=script&lang=js":
/*!*****************************************************************************************************************************************************************************!*\
  !*** ./node_modules/babel-loader/lib??ref--4-0!./node_modules/vue-loader/lib??vue-loader-options!./resources/js/components/salev2/OrderManager.vue?vue&type=script&lang=js ***!
  \*****************************************************************************************************************************************************************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
function _typeof(o) { "@babel/helpers - typeof"; return _typeof = "function" == typeof Symbol && "symbol" == typeof Symbol.iterator ? function (o) { return typeof o; } : function (o) { return o && "function" == typeof Symbol && o.constructor === Symbol && o !== Symbol.prototype ? "symbol" : typeof o; }, _typeof(o); }
function _defineProperty(e, r, t) { return (r = _toPropertyKey(r)) in e ? Object.defineProperty(e, r, { value: t, enumerable: !0, configurable: !0, writable: !0 }) : e[r] = t, e; }
function _toPropertyKey(t) { var i = _toPrimitive(t, "string"); return "symbol" == _typeof(i) ? i : i + ""; }
function _toPrimitive(t, r) { if ("object" != _typeof(t) || !t) return t; var e = t[Symbol.toPrimitive]; if (void 0 !== e) { var i = e.call(t, r || "default"); if ("object" != _typeof(i)) return i; throw new TypeError("@@toPrimitive must return a primitive value."); } return ("string" === r ? String : Number)(t); }
function _regenerator() { /*! regenerator-runtime -- Copyright (c) 2014-present, Facebook, Inc. -- license (MIT): https://github.com/babel/babel/blob/main/packages/babel-helpers/LICENSE */ var e, t, r = "function" == typeof Symbol ? Symbol : {}, n = r.iterator || "@@iterator", o = r.toStringTag || "@@toStringTag"; function i(r, n, o, i) { var c = n && n.prototype instanceof Generator ? n : Generator, u = Object.create(c.prototype); return _regeneratorDefine2(u, "_invoke", function (r, n, o) { var i, c, u, f = 0, p = o || [], y = !1, G = { p: 0, n: 0, v: e, a: d, f: d.bind(e, 4), d: function d(t, r) { return i = t, c = 0, u = e, G.n = r, a; } }; function d(r, n) { for (c = r, u = n, t = 0; !y && f && !o && t < p.length; t++) { var o, i = p[t], d = G.p, l = i[2]; r > 3 ? (o = l === n) && (u = i[(c = i[4]) ? 5 : (c = 3, 3)], i[4] = i[5] = e) : i[0] <= d && ((o = r < 2 && d < i[1]) ? (c = 0, G.v = n, G.n = i[1]) : d < l && (o = r < 3 || i[0] > n || n > l) && (i[4] = r, i[5] = n, G.n = l, c = 0)); } if (o || r > 1) return a; throw y = !0, n; } return function (o, p, l) { if (f > 1) throw TypeError("Generator is already running"); for (y && 1 === p && d(p, l), c = p, u = l; (t = c < 2 ? e : u) || !y;) { i || (c ? c < 3 ? (c > 1 && (G.n = -1), d(c, u)) : G.n = u : G.v = u); try { if (f = 2, i) { if (c || (o = "next"), t = i[o]) { if (!(t = t.call(i, u))) throw TypeError("iterator result is not an object"); if (!t.done) return t; u = t.value, c < 2 && (c = 0); } else 1 === c && (t = i["return"]) && t.call(i), c < 2 && (u = TypeError("The iterator does not provide a '" + o + "' method"), c = 1); i = e; } else if ((t = (y = G.n < 0) ? u : r.call(n, G)) !== a) break; } catch (t) { i = e, c = 1, u = t; } finally { f = 1; } } return { value: t, done: y }; }; }(r, o, i), !0), u; } var a = {}; function Generator() {} function GeneratorFunction() {} function GeneratorFunctionPrototype() {} t = Object.getPrototypeOf; var c = [][n] ? t(t([][n]())) : (_regeneratorDefine2(t = {}, n, function () { return this; }), t), u = GeneratorFunctionPrototype.prototype = Generator.prototype = Object.create(c); function f(e) { return Object.setPrototypeOf ? Object.setPrototypeOf(e, GeneratorFunctionPrototype) : (e.__proto__ = GeneratorFunctionPrototype, _regeneratorDefine2(e, o, "GeneratorFunction")), e.prototype = Object.create(u), e; } return GeneratorFunction.prototype = GeneratorFunctionPrototype, _regeneratorDefine2(u, "constructor", GeneratorFunctionPrototype), _regeneratorDefine2(GeneratorFunctionPrototype, "constructor", GeneratorFunction), GeneratorFunction.displayName = "GeneratorFunction", _regeneratorDefine2(GeneratorFunctionPrototype, o, "GeneratorFunction"), _regeneratorDefine2(u), _regeneratorDefine2(u, o, "Generator"), _regeneratorDefine2(u, n, function () { return this; }), _regeneratorDefine2(u, "toString", function () { return "[object Generator]"; }), (_regenerator = function _regenerator() { return { w: i, m: f }; })(); }
function _regeneratorDefine2(e, r, n, t) { var i = Object.defineProperty; try { i({}, "", {}); } catch (e) { i = 0; } _regeneratorDefine2 = function _regeneratorDefine(e, r, n, t) { function o(r, n) { _regeneratorDefine2(e, r, function (e) { return this._invoke(r, n, e); }); } r ? i ? i(e, r, { value: n, enumerable: !t, configurable: !t, writable: !t }) : e[r] = n : (o("next", 0), o("throw", 1), o("return", 2)); }, _regeneratorDefine2(e, r, n, t); }
function asyncGeneratorStep(n, t, e, r, o, a, c) { try { var i = n[a](c), u = i.value; } catch (n) { return void e(n); } i.done ? t(u) : Promise.resolve(u).then(r, o); }
function _asyncToGenerator(n) { return function () { var t = this, e = arguments; return new Promise(function (r, o) { var a = n.apply(t, e); function _next(n) { asyncGeneratorStep(a, r, o, _next, _throw, "next", n); } function _throw(n) { asyncGeneratorStep(a, r, o, _next, _throw, "throw", n); } _next(void 0); }); }; }
/* harmony default export */ __webpack_exports__["default"] = (_defineProperty({
  name: 'OrderManager',
  data: function data() {
    return {
      activeTab: 'NEW',
      viewMode: 'grid',
      orders: [],
      loading: false,
      total: 0,
      page: 1,
      search: '',
      currentCapture: null,
      captureVisible: false,
      captureImg: '',
      sizeMapping: {
        'VAI': 'V',
        'NGỰC': 'N',
        'EO': 'E',
        'DÀI ÁO': 'D',
        'DÀI ÁO SAU': 'DS',
        'DÀI TAY': 'DT',
        'EO QUẦN': 'E',
        'ĐÁY': 'ĐA',
        'ĐÙI': 'Đ',
        'DÀI QUẦN': 'D',
        'ỐNG': 'Ô',
        'DƯ LAI': 'L'
      }
    };
  },
  mounted: function mounted() {
    this.fetchOrders();
  },
  methods: {
    formatPrice: function formatPrice(val) {
      return val !== undefined ? new Intl.NumberFormat('vi-VN', {
        style: 'currency',
        currency: 'VND',
        maximumFractionDigits: 0
      }).format(val) : '0đ';
    },
    getStatusType: function getStatusType(s) {
      return {
        'NEW': 'info',
        'COMPLETED': 'success',
        'A SHIP NOW': 'primary',
        'ON HOLD': 'warning'
      }[s] || 'info';
    },
    formatSizes: function formatSizes(product) {
      var _this = this;
      if (!product.sizes) return [];
      var results = [];
      var duLai = product.sizes.find(function (s) {
        return s.name === 'DƯ LAI';
      });
      product.sizes.forEach(function (size) {
        if (size.name === 'DƯ LAI') return;
        var label = _this.sizeMapping[size.name] || size.name;
        var value = Math.floor(size.value);
        if (size.name === 'DÀI QUẦN' && duLai) {
          results.push("".concat(label).concat(value, "(+").concat(Math.floor(duLai.value), ")"));
        } else {
          results.push("".concat(label).concat(value));
        }
      });
      return results;
    },
    fetchOrders: function fetchOrders() {
      var _arguments = arguments,
        _this2 = this;
      return _asyncToGenerator(/*#__PURE__*/_regenerator().m(function _callee() {
        var page, params, res, _t;
        return _regenerator().w(function (_context) {
          while (1) switch (_context.p = _context.n) {
            case 0:
              page = _arguments.length > 0 && _arguments[0] !== undefined ? _arguments[0] : 1;
              _this2.page = page;
              _this2.loading = true;
              _context.p = 1;
              params = {
                page: page,
                search: _this2.search
              };
              if (_this2.activeTab === 'NEW') params.status = 'NEW';
              _context.n = 2;
              return axios.get('/salev2/api/orders', {
                params: params
              });
            case 2:
              res = _context.v;
              _this2.orders = res.data.data;
              _this2.total = res.data.total;
              _context.n = 4;
              break;
            case 3:
              _context.p = 3;
              _t = _context.v;
              console.error(_t);
            case 4:
              _context.p = 4;
              _this2.loading = false;
              return _context.f(4);
            case 5:
              return _context.a(2);
          }
        }, _callee, null, [[1, 3, 4, 5]]);
      }))();
    },
    handleTabClick: function handleTabClick() {
      this.fetchOrders(1);
    },
    debouncedFetch: _.debounce(function () {
      this.fetchOrders(1);
    }, 500),
    handlePageChange: function handlePageChange(p) {
      this.fetchOrders(p);
    },
    captureOrder: function captureOrder(order) {
      var _this3 = this;
      return _asyncToGenerator(/*#__PURE__*/_regenerator().m(function _callee3() {
        return _regenerator().w(function (_context3) {
          while (1) switch (_context3.n) {
            case 0:
              _this3.currentCapture = order;
              _this3.$nextTick(/*#__PURE__*/_asyncToGenerator(/*#__PURE__*/_regenerator().m(function _callee2() {
                var canvas;
                return _regenerator().w(function (_context2) {
                  while (1) switch (_context2.n) {
                    case 0:
                      _context2.n = 1;
                      return html2canvas(_this3.$refs.receiptTemplate, {
                        useCORS: true,
                        scale: 2
                      });
                    case 1:
                      canvas = _context2.v;
                      _this3.captureImg = canvas.toDataURL('image/jpeg', 0.9);
                      _this3.captureVisible = true;
                    case 2:
                      return _context2.a(2);
                  }
                }, _callee2);
              })));
            case 1:
              return _context3.a(2);
          }
        }, _callee3);
      }))();
    },
    downloadCapture: function downloadCapture() {
      var link = document.createElement('a');
      link.download = "Order_".concat(this.currentCapture.name, ".jpg");
      link.href = this.captureImg;
      link.click();
    },
    confirmDelete: function confirmDelete(id) {
      var _this4 = this;
      return _asyncToGenerator(/*#__PURE__*/_regenerator().m(function _callee4() {
        var _t2;
        return _regenerator().w(function (_context4) {
          while (1) switch (_context4.p = _context4.n) {
            case 0:
              if (!confirm('Delete this order?')) {
                _context4.n = 4;
                break;
              }
              _context4.p = 1;
              _context4.n = 2;
              return axios.post("/salev2/api/order/".concat(id, "/delete"));
            case 2:
              _this4.fetchOrders(_this4.page);
              _context4.n = 4;
              break;
            case 3:
              _context4.p = 3;
              _t2 = _context4.v;
            case 4:
              return _context4.a(2);
          }
        }, _callee4, null, [[1, 3]]);
      }))();
    },
    completeOrder: function completeOrder(id) {
      var _this5 = this;
      return _asyncToGenerator(/*#__PURE__*/_regenerator().m(function _callee5() {
        var _t3;
        return _regenerator().w(function (_context5) {
          while (1) switch (_context5.p = _context5.n) {
            case 0:
              _context5.p = 0;
              _context5.n = 1;
              return axios.post("/salev2/api/order/".concat(id), {
                status: 'COMPLETED'
              });
            case 1:
              _this5.fetchOrders(_this5.page);
              _context5.n = 3;
              break;
            case 2:
              _context5.p = 2;
              _t3 = _context5.v;
            case 3:
              return _context5.a(2);
          }
        }, _callee5, null, [[0, 2]]);
      }))();
    },
    updateStatus: function updateStatus(order, status) {
      var _this6 = this;
      return _asyncToGenerator(/*#__PURE__*/_regenerator().m(function _callee6() {
        var _t4;
        return _regenerator().w(function (_context6) {
          while (1) switch (_context6.p = _context6.n) {
            case 0:
              _context6.p = 0;
              _context6.n = 1;
              return axios.post("/salev2/api/order/".concat(order.id), {
                status: status
              });
            case 1:
              _this6.$message.success("Order marked as ".concat(status));
              _this6.fetchOrders(_this6.page);
              _context6.n = 3;
              break;
            case 2:
              _context6.p = 2;
              _t4 = _context6.v;
              _this6.$message.error('Update failed');
            case 3:
              return _context6.a(2);
          }
        }, _callee6, null, [[0, 2]]);
      }))();
    },
    viewCustomerHistory: function viewCustomerHistory(order) {
      // Redirect to Orders page filtered by this phone
      if (this.$route.path === '/salev2/orders') {
        this.search = order.phone;
        this.fetchOrders(1);
      } else {
        this.$router.push({
          path: '/salev2/orders',
          query: {
            search: order.phone
          }
        });
      }
    }
  }
}, "mounted", function mounted() {
  if (this.$route.query.search) {
    this.search = this.$route.query.search;
  }
  this.fetchOrders();
}));

/***/ }),

/***/ "./node_modules/babel-loader/lib/index.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/ProductCreate.vue?vue&type=script&lang=js":
/*!******************************************************************************************************************************************************************************!*\
  !*** ./node_modules/babel-loader/lib??ref--4-0!./node_modules/vue-loader/lib??vue-loader-options!./resources/js/components/salev2/ProductCreate.vue?vue&type=script&lang=js ***!
  \******************************************************************************************************************************************************************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
function _regenerator() { /*! regenerator-runtime -- Copyright (c) 2014-present, Facebook, Inc. -- license (MIT): https://github.com/babel/babel/blob/main/packages/babel-helpers/LICENSE */ var e, t, r = "function" == typeof Symbol ? Symbol : {}, n = r.iterator || "@@iterator", o = r.toStringTag || "@@toStringTag"; function i(r, n, o, i) { var c = n && n.prototype instanceof Generator ? n : Generator, u = Object.create(c.prototype); return _regeneratorDefine2(u, "_invoke", function (r, n, o) { var i, c, u, f = 0, p = o || [], y = !1, G = { p: 0, n: 0, v: e, a: d, f: d.bind(e, 4), d: function d(t, r) { return i = t, c = 0, u = e, G.n = r, a; } }; function d(r, n) { for (c = r, u = n, t = 0; !y && f && !o && t < p.length; t++) { var o, i = p[t], d = G.p, l = i[2]; r > 3 ? (o = l === n) && (u = i[(c = i[4]) ? 5 : (c = 3, 3)], i[4] = i[5] = e) : i[0] <= d && ((o = r < 2 && d < i[1]) ? (c = 0, G.v = n, G.n = i[1]) : d < l && (o = r < 3 || i[0] > n || n > l) && (i[4] = r, i[5] = n, G.n = l, c = 0)); } if (o || r > 1) return a; throw y = !0, n; } return function (o, p, l) { if (f > 1) throw TypeError("Generator is already running"); for (y && 1 === p && d(p, l), c = p, u = l; (t = c < 2 ? e : u) || !y;) { i || (c ? c < 3 ? (c > 1 && (G.n = -1), d(c, u)) : G.n = u : G.v = u); try { if (f = 2, i) { if (c || (o = "next"), t = i[o]) { if (!(t = t.call(i, u))) throw TypeError("iterator result is not an object"); if (!t.done) return t; u = t.value, c < 2 && (c = 0); } else 1 === c && (t = i["return"]) && t.call(i), c < 2 && (u = TypeError("The iterator does not provide a '" + o + "' method"), c = 1); i = e; } else if ((t = (y = G.n < 0) ? u : r.call(n, G)) !== a) break; } catch (t) { i = e, c = 1, u = t; } finally { f = 1; } } return { value: t, done: y }; }; }(r, o, i), !0), u; } var a = {}; function Generator() {} function GeneratorFunction() {} function GeneratorFunctionPrototype() {} t = Object.getPrototypeOf; var c = [][n] ? t(t([][n]())) : (_regeneratorDefine2(t = {}, n, function () { return this; }), t), u = GeneratorFunctionPrototype.prototype = Generator.prototype = Object.create(c); function f(e) { return Object.setPrototypeOf ? Object.setPrototypeOf(e, GeneratorFunctionPrototype) : (e.__proto__ = GeneratorFunctionPrototype, _regeneratorDefine2(e, o, "GeneratorFunction")), e.prototype = Object.create(u), e; } return GeneratorFunction.prototype = GeneratorFunctionPrototype, _regeneratorDefine2(u, "constructor", GeneratorFunctionPrototype), _regeneratorDefine2(GeneratorFunctionPrototype, "constructor", GeneratorFunction), GeneratorFunction.displayName = "GeneratorFunction", _regeneratorDefine2(GeneratorFunctionPrototype, o, "GeneratorFunction"), _regeneratorDefine2(u), _regeneratorDefine2(u, o, "Generator"), _regeneratorDefine2(u, n, function () { return this; }), _regeneratorDefine2(u, "toString", function () { return "[object Generator]"; }), (_regenerator = function _regenerator() { return { w: i, m: f }; })(); }
function _regeneratorDefine2(e, r, n, t) { var i = Object.defineProperty; try { i({}, "", {}); } catch (e) { i = 0; } _regeneratorDefine2 = function _regeneratorDefine(e, r, n, t) { function o(r, n) { _regeneratorDefine2(e, r, function (e) { return this._invoke(r, n, e); }); } r ? i ? i(e, r, { value: n, enumerable: !t, configurable: !t, writable: !t }) : e[r] = n : (o("next", 0), o("throw", 1), o("return", 2)); }, _regeneratorDefine2(e, r, n, t); }
function asyncGeneratorStep(n, t, e, r, o, a, c) { try { var i = n[a](c), u = i.value; } catch (n) { return void e(n); } i.done ? t(u) : Promise.resolve(u).then(r, o); }
function _asyncToGenerator(n) { return function () { var t = this, e = arguments; return new Promise(function (r, o) { var a = n.apply(t, e); function _next(n) { asyncGeneratorStep(a, r, o, _next, _throw, "next", n); } function _throw(n) { asyncGeneratorStep(a, r, o, _next, _throw, "throw", n); } _next(void 0); }); }; }
/* harmony default export */ __webpack_exports__["default"] = ({
  name: 'ProductCreate',
  data: function data() {
    return {
      form: {
        type: '',
        price: 0,
        description: '',
        sizes: [],
        image: null
      },
      quickInput: '',
      previewUrl: null,
      loading: false,
      types: ['TROUSERS', 'JACKET', 'SHIRT', 'BLAZER', 'TIE', 'GILE', 'BELT', 'POLO SHIRT', 'HAT', 'SUIT'],
      // Copied from legacy JS logic and SaleConfig
      mappings: {
        'BLAZER': {
          labels: 'Vai - Ngực - Dài trước - Dài tay',
          indices: [0, 1, 3, 5]
        },
        'JACKET': {
          labels: 'Vai - Ngực - Dài - Dài tay',
          indices: [0, 1, 3, 5]
        },
        'TROUSERS': {
          labels: 'Eo - Đáy - Đùi - Dài quần - Ống - Dư lai',
          indices: [6, 7, 8, 9, 10, 11]
        },
        'SHIRT': {
          labels: 'Vai - Ngực - Dài - Dài tay',
          indices: [0, 1, 3, 5]
        },
        'GILE': {
          labels: 'Vai - Ngực - Dài - Dài tay',
          indices: [0, 1, 3, 5]
        },
        'POLO SHIRT': {
          labels: 'Vai - Ngực - Dài - Dài tay',
          indices: [0, 1, 3, 5]
        },
        'SUIT': {
          labels: 'Vai - Ngực - Dài - Dài tay',
          indices: [0, 1, 3, 5]
        }
      },
      // The global list of sizes from legacy Contract
      allSizes: ['VAI', 'NGỰC', 'EO', 'DÀI ÁO', 'DÀI ÁO SAU', 'DÀI TAY',
      // Top (Indices 0-5)
      'EO QUẦN', 'ĐÁY', 'ĐÙI', 'DÀI QUẦN', 'ỐNG', 'DƯ LAI' // Bottom (Indices 6-11)
      ]
    };
  },
  computed: {
    sizeConfig: function sizeConfig() {
      return this.mappings[this.form.type] || {
        labels: '',
        indices: []
      };
    },
    activeSizes: function activeSizes() {
      // We need to determine which sizes to show based on type. 
      // Legacy logic hides/shows sections: AO (< 6) vs QUAN (>= 6).
      if (['TROUSERS', 'SHORT'].includes(this.form.type)) {
        return this.allSizes.slice(6); // The bottom half
      } else {
        return this.allSizes.slice(0, 6); // The top half
      }
    }
  },
  methods: {
    handleTypeChange: function handleTypeChange() {
      // Reset sizes array to match the global structure but we only fill what's needed
      // Actually, we should initialize form.sizes with {name, value} for ALL sizes to match legacy payload structure?
      // Legacy payload: sizes[0][name]=Vai, sizes[0][value]=...
      // So we should maintain the full array of 12 items.
      this.form.sizes = this.allSizes.map(function (name) {
        return {
          name: name,
          value: ''
        };
      });
      this.quickInput = '';
    },
    parseQuickInput: function parseQuickInput(val) {
      var _this = this;
      if (!val) return;
      var numbers = val.trim().split(/\s+/);
      var indices = this.sizeConfig.indices;
      if (!indices || !indices.length) return;
      indices.forEach(function (globalIdx, i) {
        if (numbers[i] && _this.form.sizes[globalIdx]) {
          _this.form.sizes[globalIdx].value = numbers[i];
        }
      });
    },
    handleFileChange: function handleFileChange(e) {
      var file = e.target.files[0];
      if (file) {
        this.form.image = file;
        this.previewUrl = URL.createObjectURL(file);
      }
    },
    submitByAjax: function submitByAjax() {
      var _this2 = this;
      return _asyncToGenerator(/*#__PURE__*/_regenerator().m(function _callee() {
        var formData, _t;
        return _regenerator().w(function (_context) {
          while (1) switch (_context.p = _context.n) {
            case 0:
              if (!(!_this2.form.type || !_this2.form.image)) {
                _context.n = 1;
                break;
              }
              return _context.a(2, _this2.$message.error('Type and Image are required'));
            case 1:
              _this2.loading = true;
              _context.p = 2;
              formData = new FormData();
              formData.append('type', _this2.form.type);
              formData.append('price', _this2.form.price);
              formData.append('description', _this2.form.description || '');
              formData.append('banner', _this2.form.image);

              // Append sizes as array of objects
              _this2.form.sizes.forEach(function (s, i) {
                formData.append("sizes[".concat(i, "][name]"), s.name);
                formData.append("sizes[".concat(i, "][value]"), s.value);
              });
              _context.n = 3;
              return axios.post('/salev2/api/product/store', formData, {
                headers: {
                  'Content-Type': 'multipart/form-data'
                }
              });
            case 3:
              _this2.$message.success('Product Created Successfully!');
              _this2.$router.push('/salev2/products');
              _context.n = 5;
              break;
            case 4:
              _context.p = 4;
              _t = _context.v;
              _this2.$message.error('Failed to create product');
            case 5:
              _context.p = 5;
              _this2.loading = false;
              return _context.f(5);
            case 6:
              return _context.a(2);
          }
        }, _callee, null, [[2, 4, 5, 6]]);
      }))();
    }
  },
  mounted: function mounted() {
    this.handleTypeChange(); // Init empty sizes
  }
});

/***/ }),

/***/ "./node_modules/babel-loader/lib/index.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/ProductEdit.vue?vue&type=script&lang=js":
/*!****************************************************************************************************************************************************************************!*\
  !*** ./node_modules/babel-loader/lib??ref--4-0!./node_modules/vue-loader/lib??vue-loader-options!./resources/js/components/salev2/ProductEdit.vue?vue&type=script&lang=js ***!
  \****************************************************************************************************************************************************************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
function _regenerator() { /*! regenerator-runtime -- Copyright (c) 2014-present, Facebook, Inc. -- license (MIT): https://github.com/babel/babel/blob/main/packages/babel-helpers/LICENSE */ var e, t, r = "function" == typeof Symbol ? Symbol : {}, n = r.iterator || "@@iterator", o = r.toStringTag || "@@toStringTag"; function i(r, n, o, i) { var c = n && n.prototype instanceof Generator ? n : Generator, u = Object.create(c.prototype); return _regeneratorDefine2(u, "_invoke", function (r, n, o) { var i, c, u, f = 0, p = o || [], y = !1, G = { p: 0, n: 0, v: e, a: d, f: d.bind(e, 4), d: function d(t, r) { return i = t, c = 0, u = e, G.n = r, a; } }; function d(r, n) { for (c = r, u = n, t = 0; !y && f && !o && t < p.length; t++) { var o, i = p[t], d = G.p, l = i[2]; r > 3 ? (o = l === n) && (u = i[(c = i[4]) ? 5 : (c = 3, 3)], i[4] = i[5] = e) : i[0] <= d && ((o = r < 2 && d < i[1]) ? (c = 0, G.v = n, G.n = i[1]) : d < l && (o = r < 3 || i[0] > n || n > l) && (i[4] = r, i[5] = n, G.n = l, c = 0)); } if (o || r > 1) return a; throw y = !0, n; } return function (o, p, l) { if (f > 1) throw TypeError("Generator is already running"); for (y && 1 === p && d(p, l), c = p, u = l; (t = c < 2 ? e : u) || !y;) { i || (c ? c < 3 ? (c > 1 && (G.n = -1), d(c, u)) : G.n = u : G.v = u); try { if (f = 2, i) { if (c || (o = "next"), t = i[o]) { if (!(t = t.call(i, u))) throw TypeError("iterator result is not an object"); if (!t.done) return t; u = t.value, c < 2 && (c = 0); } else 1 === c && (t = i["return"]) && t.call(i), c < 2 && (u = TypeError("The iterator does not provide a '" + o + "' method"), c = 1); i = e; } else if ((t = (y = G.n < 0) ? u : r.call(n, G)) !== a) break; } catch (t) { i = e, c = 1, u = t; } finally { f = 1; } } return { value: t, done: y }; }; }(r, o, i), !0), u; } var a = {}; function Generator() {} function GeneratorFunction() {} function GeneratorFunctionPrototype() {} t = Object.getPrototypeOf; var c = [][n] ? t(t([][n]())) : (_regeneratorDefine2(t = {}, n, function () { return this; }), t), u = GeneratorFunctionPrototype.prototype = Generator.prototype = Object.create(c); function f(e) { return Object.setPrototypeOf ? Object.setPrototypeOf(e, GeneratorFunctionPrototype) : (e.__proto__ = GeneratorFunctionPrototype, _regeneratorDefine2(e, o, "GeneratorFunction")), e.prototype = Object.create(u), e; } return GeneratorFunction.prototype = GeneratorFunctionPrototype, _regeneratorDefine2(u, "constructor", GeneratorFunctionPrototype), _regeneratorDefine2(GeneratorFunctionPrototype, "constructor", GeneratorFunction), GeneratorFunction.displayName = "GeneratorFunction", _regeneratorDefine2(GeneratorFunctionPrototype, o, "GeneratorFunction"), _regeneratorDefine2(u), _regeneratorDefine2(u, o, "Generator"), _regeneratorDefine2(u, n, function () { return this; }), _regeneratorDefine2(u, "toString", function () { return "[object Generator]"; }), (_regenerator = function _regenerator() { return { w: i, m: f }; })(); }
function _regeneratorDefine2(e, r, n, t) { var i = Object.defineProperty; try { i({}, "", {}); } catch (e) { i = 0; } _regeneratorDefine2 = function _regeneratorDefine(e, r, n, t) { function o(r, n) { _regeneratorDefine2(e, r, function (e) { return this._invoke(r, n, e); }); } r ? i ? i(e, r, { value: n, enumerable: !t, configurable: !t, writable: !t }) : e[r] = n : (o("next", 0), o("throw", 1), o("return", 2)); }, _regeneratorDefine2(e, r, n, t); }
function asyncGeneratorStep(n, t, e, r, o, a, c) { try { var i = n[a](c), u = i.value; } catch (n) { return void e(n); } i.done ? t(u) : Promise.resolve(u).then(r, o); }
function _asyncToGenerator(n) { return function () { var t = this, e = arguments; return new Promise(function (r, o) { var a = n.apply(t, e); function _next(n) { asyncGeneratorStep(a, r, o, _next, _throw, "next", n); } function _throw(n) { asyncGeneratorStep(a, r, o, _next, _throw, "throw", n); } _next(void 0); }); }; }
/* harmony default export */ __webpack_exports__["default"] = ({
  name: 'ProductEdit',
  data: function data() {
    return {
      loading: false,
      submitting: false,
      form: {
        id: null,
        type: '',
        price: 0,
        description: '',
        image: null,
        sizes: []
      },
      imagePreview: '',
      quickInput: '',
      productTypes: ['TROUSERS', 'JACKET', 'SHIRT', 'BLAZER', 'TIE', 'GILE', 'BELT', 'POLO SHIRT', 'HAT', 'SUIT'],
      allSizes: ['VAI', 'NGỰC', 'EO', 'DÀI ÁO', 'DÀI ÁO SAU', 'DÀI TAY', 'EO QUẦN', 'ĐÁY', 'ĐÙI', 'DÀI QUẦN', 'ỐNG', 'DƯ LAI'],
      mappings: {
        'BLAZER': {
          labels: 'Vai - Ngực - Dài trước - Dài tay',
          indices: [0, 1, 3, 5]
        },
        'JACKET': {
          labels: 'Vai - Ngực - Dài - Dài tay',
          indices: [0, 1, 3, 5]
        },
        'TROUSERS': {
          labels: 'Eo - Đáy - Đùi - Dài quần - Ống - Dư lai',
          indices: [6, 7, 8, 9, 10, 11]
        },
        'SHIRT': {
          labels: 'Vai - Ngực - Dài - Dài tay',
          indices: [0, 1, 3, 5]
        },
        'GILE': {
          labels: 'Vai - Ngực - Dài - Dài tay',
          indices: [0, 1, 3, 5]
        },
        'POLO SHIRT': {
          labels: 'Vai - Ngực - Dài - Dài tay',
          indices: [0, 1, 3, 5]
        },
        'SUIT': {
          labels: 'Vai - Ngực - Dài - Dài tay',
          indices: [0, 1, 3, 5]
        }
      }
    };
  },
  computed: {
    activeMapping: function activeMapping() {
      return this.mappings[this.form.type] || null;
    }
  },
  mounted: function mounted() {
    this.initializeSizes();
    this.fetchProduct();
  },
  methods: {
    initializeSizes: function initializeSizes() {
      this.form.sizes = this.allSizes.map(function (name) {
        return {
          name: name,
          value: ''
        };
      });
    },
    fetchProduct: function fetchProduct() {
      var _this = this;
      return _asyncToGenerator(/*#__PURE__*/_regenerator().m(function _callee() {
        var id, res, product, _t;
        return _regenerator().w(function (_context) {
          while (1) switch (_context.p = _context.n) {
            case 0:
              _this.loading = true;
              _context.p = 1;
              id = _this.$route.params.id;
              _context.n = 2;
              return axios.get("/salev2/api/product/".concat(id));
            case 2:
              res = _context.v;
              product = res.data;
              _this.form.id = product.id;
              _this.form.type = product.type;
              _this.form.price = product.price;
              _this.form.description = product.description;
              _this.imagePreview = product.path_thumb;

              // Populate sizes
              if (product.sizes) {
                product.sizes.forEach(function (s) {
                  var idx = _this.allSizes.indexOf(s.name);
                  if (idx !== -1) _this.form.sizes[idx].value = s.value;
                });
              }
              _context.n = 4;
              break;
            case 3:
              _context.p = 3;
              _t = _context.v;
              _this.$message.error('Failed to load product');
              _this.$router.push('/salev2/products');
            case 4:
              _context.p = 4;
              _this.loading = false;
              return _context.f(4);
            case 5:
              return _context.a(2);
          }
        }, _callee, null, [[1, 3, 4, 5]]);
      }))();
    },
    handleFileChange: function handleFileChange(e) {
      var file = e.target.files[0];
      if (file) {
        this.form.image = file;
        this.imagePreview = URL.createObjectURL(file);
      }
    },
    isSizeRelevant: function isSizeRelevant(idx) {
      if (!this.form.type) return false;
      // Simple logic: Trousers/Shorts use bottom sizes (idx >= 6), others use top sizes (idx < 6)
      var bottomTypes = ['TROUSERS', 'SHORT'];
      var isBottom = bottomTypes.includes(this.form.type);
      return isBottom ? idx >= 6 : idx < 6;
    },
    parseQuickInput: function parseQuickInput() {
      var _this2 = this;
      if (!this.quickInput) return;
      var nums = this.quickInput.trim().split(/\s+/);
      var config = this.activeMapping;
      if (config && config.indices) {
        config.indices.forEach(function (sizeIdx, i) {
          if (nums[i]) {
            _this2.form.sizes[sizeIdx].value = nums[i];
          }
        });
      }
    },
    submit: function submit() {
      var _this3 = this;
      return _asyncToGenerator(/*#__PURE__*/_regenerator().m(function _callee2() {
        var formData, _t2;
        return _regenerator().w(function (_context2) {
          while (1) switch (_context2.p = _context2.n) {
            case 0:
              _this3.submitting = true;
              _context2.p = 1;
              formData = new FormData();
              formData.append('type', _this3.form.type);
              formData.append('price', _this3.form.price);
              formData.append('description', _this3.form.description || '');
              if (_this3.form.image) {
                formData.append('banner', _this3.form.image);
              }
              _this3.form.sizes.forEach(function (s, i) {
                formData.append("sizes[".concat(i, "][name]"), s.name);
                formData.append("sizes[".concat(i, "][value]"), s.value);
              });
              _context2.n = 2;
              return axios.post("/salev2/api/product/".concat(_this3.form.id, "/update"), formData, {
                headers: {
                  'Content-Type': 'multipart/form-data'
                }
              });
            case 2:
              _this3.$message.success('Product updated successfully');
              _this3.$router.push('/salev2/products');
              _context2.n = 4;
              break;
            case 3:
              _context2.p = 3;
              _t2 = _context2.v;
              console.error(_t2);
              _this3.$message.error('Failed to update product');
            case 4:
              _context2.p = 4;
              _this3.submitting = false;
              return _context2.f(4);
            case 5:
              return _context2.a(2);
          }
        }, _callee2, null, [[1, 3, 4, 5]]);
      }))();
    }
  }
});

/***/ }),

/***/ "./node_modules/babel-loader/lib/index.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/ProductVault.vue?vue&type=script&lang=js":
/*!*****************************************************************************************************************************************************************************!*\
  !*** ./node_modules/babel-loader/lib??ref--4-0!./node_modules/vue-loader/lib??vue-loader-options!./resources/js/components/salev2/ProductVault.vue?vue&type=script&lang=js ***!
  \*****************************************************************************************************************************************************************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
function _typeof(o) { "@babel/helpers - typeof"; return _typeof = "function" == typeof Symbol && "symbol" == typeof Symbol.iterator ? function (o) { return typeof o; } : function (o) { return o && "function" == typeof Symbol && o.constructor === Symbol && o !== Symbol.prototype ? "symbol" : typeof o; }, _typeof(o); }
function _defineProperty(e, r, t) { return (r = _toPropertyKey(r)) in e ? Object.defineProperty(e, r, { value: t, enumerable: !0, configurable: !0, writable: !0 }) : e[r] = t, e; }
function _toPropertyKey(t) { var i = _toPrimitive(t, "string"); return "symbol" == _typeof(i) ? i : i + ""; }
function _toPrimitive(t, r) { if ("object" != _typeof(t) || !t) return t; var e = t[Symbol.toPrimitive]; if (void 0 !== e) { var i = e.call(t, r || "default"); if ("object" != _typeof(i)) return i; throw new TypeError("@@toPrimitive must return a primitive value."); } return ("string" === r ? String : Number)(t); }
function _createForOfIteratorHelper(r, e) { var t = "undefined" != typeof Symbol && r[Symbol.iterator] || r["@@iterator"]; if (!t) { if (Array.isArray(r) || (t = _unsupportedIterableToArray(r)) || e && r && "number" == typeof r.length) { t && (r = t); var _n = 0, F = function F() {}; return { s: F, n: function n() { return _n >= r.length ? { done: !0 } : { done: !1, value: r[_n++] }; }, e: function e(r) { throw r; }, f: F }; } throw new TypeError("Invalid attempt to iterate non-iterable instance.\nIn order to be iterable, non-array objects must have a [Symbol.iterator]() method."); } var o, a = !0, u = !1; return { s: function s() { t = t.call(r); }, n: function n() { var r = t.next(); return a = r.done, r; }, e: function e(r) { u = !0, o = r; }, f: function f() { try { a || null == t["return"] || t["return"](); } finally { if (u) throw o; } } }; }
function _unsupportedIterableToArray(r, a) { if (r) { if ("string" == typeof r) return _arrayLikeToArray(r, a); var t = {}.toString.call(r).slice(8, -1); return "Object" === t && r.constructor && (t = r.constructor.name), "Map" === t || "Set" === t ? Array.from(r) : "Arguments" === t || /^(?:Ui|I)nt(?:8|16|32)(?:Clamped)?Array$/.test(t) ? _arrayLikeToArray(r, a) : void 0; } }
function _arrayLikeToArray(r, a) { (null == a || a > r.length) && (a = r.length); for (var e = 0, n = Array(a); e < a; e++) n[e] = r[e]; return n; }
function _regenerator() { /*! regenerator-runtime -- Copyright (c) 2014-present, Facebook, Inc. -- license (MIT): https://github.com/babel/babel/blob/main/packages/babel-helpers/LICENSE */ var e, t, r = "function" == typeof Symbol ? Symbol : {}, n = r.iterator || "@@iterator", o = r.toStringTag || "@@toStringTag"; function i(r, n, o, i) { var c = n && n.prototype instanceof Generator ? n : Generator, u = Object.create(c.prototype); return _regeneratorDefine2(u, "_invoke", function (r, n, o) { var i, c, u, f = 0, p = o || [], y = !1, G = { p: 0, n: 0, v: e, a: d, f: d.bind(e, 4), d: function d(t, r) { return i = t, c = 0, u = e, G.n = r, a; } }; function d(r, n) { for (c = r, u = n, t = 0; !y && f && !o && t < p.length; t++) { var o, i = p[t], d = G.p, l = i[2]; r > 3 ? (o = l === n) && (u = i[(c = i[4]) ? 5 : (c = 3, 3)], i[4] = i[5] = e) : i[0] <= d && ((o = r < 2 && d < i[1]) ? (c = 0, G.v = n, G.n = i[1]) : d < l && (o = r < 3 || i[0] > n || n > l) && (i[4] = r, i[5] = n, G.n = l, c = 0)); } if (o || r > 1) return a; throw y = !0, n; } return function (o, p, l) { if (f > 1) throw TypeError("Generator is already running"); for (y && 1 === p && d(p, l), c = p, u = l; (t = c < 2 ? e : u) || !y;) { i || (c ? c < 3 ? (c > 1 && (G.n = -1), d(c, u)) : G.n = u : G.v = u); try { if (f = 2, i) { if (c || (o = "next"), t = i[o]) { if (!(t = t.call(i, u))) throw TypeError("iterator result is not an object"); if (!t.done) return t; u = t.value, c < 2 && (c = 0); } else 1 === c && (t = i["return"]) && t.call(i), c < 2 && (u = TypeError("The iterator does not provide a '" + o + "' method"), c = 1); i = e; } else if ((t = (y = G.n < 0) ? u : r.call(n, G)) !== a) break; } catch (t) { i = e, c = 1, u = t; } finally { f = 1; } } return { value: t, done: y }; }; }(r, o, i), !0), u; } var a = {}; function Generator() {} function GeneratorFunction() {} function GeneratorFunctionPrototype() {} t = Object.getPrototypeOf; var c = [][n] ? t(t([][n]())) : (_regeneratorDefine2(t = {}, n, function () { return this; }), t), u = GeneratorFunctionPrototype.prototype = Generator.prototype = Object.create(c); function f(e) { return Object.setPrototypeOf ? Object.setPrototypeOf(e, GeneratorFunctionPrototype) : (e.__proto__ = GeneratorFunctionPrototype, _regeneratorDefine2(e, o, "GeneratorFunction")), e.prototype = Object.create(u), e; } return GeneratorFunction.prototype = GeneratorFunctionPrototype, _regeneratorDefine2(u, "constructor", GeneratorFunctionPrototype), _regeneratorDefine2(GeneratorFunctionPrototype, "constructor", GeneratorFunction), GeneratorFunction.displayName = "GeneratorFunction", _regeneratorDefine2(GeneratorFunctionPrototype, o, "GeneratorFunction"), _regeneratorDefine2(u), _regeneratorDefine2(u, o, "Generator"), _regeneratorDefine2(u, n, function () { return this; }), _regeneratorDefine2(u, "toString", function () { return "[object Generator]"; }), (_regenerator = function _regenerator() { return { w: i, m: f }; })(); }
function _regeneratorDefine2(e, r, n, t) { var i = Object.defineProperty; try { i({}, "", {}); } catch (e) { i = 0; } _regeneratorDefine2 = function _regeneratorDefine(e, r, n, t) { function o(r, n) { _regeneratorDefine2(e, r, function (e) { return this._invoke(r, n, e); }); } r ? i ? i(e, r, { value: n, enumerable: !t, configurable: !t, writable: !t }) : e[r] = n : (o("next", 0), o("throw", 1), o("return", 2)); }, _regeneratorDefine2(e, r, n, t); }
function asyncGeneratorStep(n, t, e, r, o, a, c) { try { var i = n[a](c), u = i.value; } catch (n) { return void e(n); } i.done ? t(u) : Promise.resolve(u).then(r, o); }
function _asyncToGenerator(n) { return function () { var t = this, e = arguments; return new Promise(function (r, o) { var a = n.apply(t, e); function _next(n) { asyncGeneratorStep(a, r, o, _next, _throw, "next", n); } function _throw(n) { asyncGeneratorStep(a, r, o, _next, _throw, "throw", n); } _next(void 0); }); }; }
/* harmony default export */ __webpack_exports__["default"] = ({
  name: 'ProductVault',
  data: function data() {
    return {
      activeTab: 'AVAILABLE',
      mobileFiltersOpen: false,
      products: [],
      loading: false,
      total: 0,
      page: 1,
      search: '',
      filterType: '',
      sizeFilters: [{
        key: '',
        from: undefined,
        to: undefined
      }],
      selectedProducts: [],
      productTypes: ['TROUSERS', 'JACKET', 'SHIRT', 'BLAZER', 'TIE', 'GILE', 'BELT', 'POLO SHIRT', 'HAT', 'SUIT'],
      sizes: ['VAI', 'NGỰC', 'EO', 'DÀI ÁO', 'DÀI ÁO SAU', 'DÀI TAY', 'EO QUẦN', 'ĐÁY', 'ĐÙI', 'DÀI QUẦN', 'ỐNG', 'DƯ LAI'],
      sizeMapping: {
        'VAI': 'V',
        'NGỰC': 'N',
        'EO': 'E',
        'DÀI ÁO': 'D',
        'DÀI ÁO SAU': 'DS',
        'DÀI TAY': 'DT',
        'EO QUẦN': 'E',
        'ĐÁY': 'ĐA',
        'ĐÙI': 'Đ',
        'DÀI QUẦN': 'D',
        'ỐNG': 'Ô',
        'DƯ LAI': 'L'
      },
      // Edit Modal Data
      editVisible: false,
      editLoading: false,
      editQuickInput: '',
      editForm: {
        id: null,
        type: '',
        price: 0,
        description: '',
        sizes: [],
        image: null,
        imagePreview: ''
      },
      mappings: {
        'BLAZER': {
          labels: 'Vai - Ngực - Dài trước - Dài tay',
          indices: [0, 1, 3, 5]
        },
        'JACKET': {
          labels: 'Vai - Ngực - Dài - Dài tay',
          indices: [0, 1, 3, 5]
        },
        'TROUSERS': {
          labels: 'Eo - Đáy - Đùi - Dài quần - Ống - Dư lai',
          indices: [6, 7, 8, 9, 10, 11]
        },
        'SHIRT': {
          labels: 'Vai - Ngực - Dài - Dài tay',
          indices: [0, 1, 3, 5]
        },
        'GILE': {
          labels: 'Vai - Ngực - Dài - Dài tay',
          indices: [0, 1, 3, 5]
        },
        'POLO SHIRT': {
          labels: 'Vai - Ngực - Dài - Dài tay',
          indices: [0, 1, 3, 5]
        },
        'SUIT': {
          labels: 'Vai - Ngực - Dài - Dài tay',
          indices: [0, 1, 3, 5]
        }
      },
      allSizes: ['VAI', 'NGỰC', 'EO', 'DÀI ÁO', 'DÀI ÁO SAU', 'DÀI TAY', 'EO QUẦN', 'ĐÁY', 'ĐÙI', 'DÀI QUẦN', 'ỐNG', 'DƯ LAI']
    };
  },
  mounted: function mounted() {
    this.fetchProducts();
    var stored = localStorage.getItem('selectedProducts');
    if (stored) this.selectedProducts = JSON.parse(stored);
  },
  methods: _defineProperty(_defineProperty(_defineProperty(_defineProperty({
    formatPrice: function formatPrice(val) {
      if (!val) return '0đ';
      return new Intl.NumberFormat('vi-VN', {
        style: 'currency',
        currency: 'VND',
        maximumFractionDigits: 0
      }).format(val);
    },
    getStatusTagType: function getStatusTagType(status) {
      var maps = {
        'AVAILABLE': 'success',
        'SOLD': 'info',
        'ON_HOLD': 'warning'
      };
      return maps[status] || 'info';
    },
    formatSizes: function formatSizes(product) {
      var _this = this;
      if (!product.sizes) return [];
      var results = [];
      var duLai = product.sizes.find(function (s) {
        return s.name === 'DƯ LAI';
      });
      product.sizes.forEach(function (size) {
        if (size.name === 'DƯ LAI') return;
        var label = _this.sizeMapping[size.name] || size.name;
        var value = Math.floor(size.value);
        if (size.name === 'DÀI QUẦN' && duLai) {
          results.push("".concat(label).concat(value, "(+").concat(Math.floor(duLai.value), ")"));
        } else {
          results.push("".concat(label).concat(value));
        }
      });
      return results;
    },
    addSizeFilter: function addSizeFilter() {
      if (this.sizeFilters.length < 3) {
        this.sizeFilters.push({
          key: '',
          from: undefined,
          to: undefined
        });
      }
    },
    removeSizeFilter: function removeSizeFilter(index) {
      this.sizeFilters.splice(index, 1);
    },
    fetchProducts: function fetchProducts() {
      var _arguments = arguments,
        _this2 = this;
      return _asyncToGenerator(/*#__PURE__*/_regenerator().m(function _callee() {
        var page, _this2$sizeFilters$, _this2$sizeFilters$2, _this2$sizeFilters$3, _this2$sizeFilters$4, _this2$sizeFilters$5, _this2$sizeFilters$6, params, res, _t;
        return _regenerator().w(function (_context) {
          while (1) switch (_context.p = _context.n) {
            case 0:
              page = _arguments.length > 0 && _arguments[0] !== undefined ? _arguments[0] : 1;
              _this2.page = page;
              _this2.loading = true;
              _context.p = 1;
              params = {
                page: page,
                search: _this2.search,
                type: _this2.filterType,
                status: _this2.activeTab === 'ALL' ? '' : _this2.activeTab,
                from_search_key: (_this2$sizeFilters$ = _this2.sizeFilters[0]) === null || _this2$sizeFilters$ === void 0 ? void 0 : _this2$sizeFilters$.key,
                from_search_value: (_this2$sizeFilters$2 = _this2.sizeFilters[0]) === null || _this2$sizeFilters$2 === void 0 ? void 0 : _this2$sizeFilters$2.from,
                to_search_value: (_this2$sizeFilters$3 = _this2.sizeFilters[0]) === null || _this2$sizeFilters$3 === void 0 ? void 0 : _this2$sizeFilters$3.to,
                from_search_key_2: (_this2$sizeFilters$4 = _this2.sizeFilters[1]) === null || _this2$sizeFilters$4 === void 0 ? void 0 : _this2$sizeFilters$4.key,
                from_search_value_2: (_this2$sizeFilters$5 = _this2.sizeFilters[1]) === null || _this2$sizeFilters$5 === void 0 ? void 0 : _this2$sizeFilters$5.from,
                to_search_value_2: (_this2$sizeFilters$6 = _this2.sizeFilters[1]) === null || _this2$sizeFilters$6 === void 0 ? void 0 : _this2$sizeFilters$6.to
              };
              _context.n = 2;
              return axios.get('/salev2/api/products', {
                params: params
              });
            case 2:
              res = _context.v;
              _this2.products = res.data.data;
              _this2.total = res.data.total;
              _context.n = 4;
              break;
            case 3:
              _context.p = 3;
              _t = _context.v;
              console.error(_t);
            case 4:
              _context.p = 4;
              _this2.loading = false;
              return _context.f(4);
            case 5:
              return _context.a(2);
          }
        }, _callee, null, [[1, 3, 4, 5]]);
      }))();
    },
    handleTabClick: function handleTabClick() {
      this.fetchProducts(1);
    },
    debouncedFetch: _.debounce(function () {
      this.fetchProducts(1);
    }, 500),
    handlePageChange: function handlePageChange(page) {
      this.fetchProducts(page);
    },
    resetFilters: function resetFilters() {
      this.search = '';
      this.filterType = '';
      this.sizeFilters = [{
        key: '',
        from: undefined,
        to: undefined
      }];
      this.fetchProducts(1);
    },
    isSelected: function isSelected(id) {
      return this.selectedProducts.some(function (p) {
        return p.id === id;
      });
    },
    toggleSelection: function toggleSelection(product) {
      var idx = this.selectedProducts.findIndex(function (p) {
        return p.id === product.id;
      });
      if (idx > -1) {
        this.selectedProducts.splice(idx, 1);
      } else {
        this.selectedProducts.push({
          id: product.id,
          img: product.path_thumb
        });
      }
      localStorage.setItem('selectedProducts', JSON.stringify(this.selectedProducts));
    },
    selectAllPage: function selectAllPage() {
      var _this3 = this;
      this.products.forEach(function (p) {
        if (!_this3.isSelected(p.id)) {
          _this3.selectedProducts.push({
            id: p.id,
            img: p.path_thumb
          });
        }
      });
      localStorage.setItem('selectedProducts', JSON.stringify(this.selectedProducts));
    },
    clearSelection: function clearSelection() {
      this.selectedProducts = [];
      localStorage.removeItem('selectedProducts');
    },
    shareSelected: function shareSelected() {
      var _this4 = this;
      return _asyncToGenerator(/*#__PURE__*/_regenerator().m(function _callee2() {
        var loading, files, _iterator, _step, item, response, blob, _t2, _t3;
        return _regenerator().w(function (_context2) {
          while (1) switch (_context2.p = _context2.n) {
            case 0:
              if (navigator.share) {
                _context2.n = 1;
                break;
              }
              _this4.$message.warning('Browser does not support native sharing');
              return _context2.a(2);
            case 1:
              loading = _this4.$loading({
                lock: true,
                text: 'Preparing images...',
                background: 'rgba(0,0,0,0.7)'
              });
              _context2.p = 2;
              files = [];
              _iterator = _createForOfIteratorHelper(_this4.selectedProducts);
              _context2.p = 3;
              _iterator.s();
            case 4:
              if ((_step = _iterator.n()).done) {
                _context2.n = 8;
                break;
              }
              item = _step.value;
              _context2.n = 5;
              return fetch(item.img);
            case 5:
              response = _context2.v;
              _context2.n = 6;
              return response.blob();
            case 6:
              blob = _context2.v;
              files.push(new File([blob], "".concat(item.id, ".jpg"), {
                type: 'image/jpeg'
              }));
            case 7:
              _context2.n = 4;
              break;
            case 8:
              _context2.n = 10;
              break;
            case 9:
              _context2.p = 9;
              _t2 = _context2.v;
              _iterator.e(_t2);
            case 10:
              _context2.p = 10;
              _iterator.f();
              return _context2.f(10);
            case 11:
              _context2.n = 12;
              return navigator.share({
                files: files,
                title: 'Vintage Always True',
                text: 'Check out these items!'
              });
            case 12:
              _context2.n = 14;
              break;
            case 13:
              _context2.p = 13;
              _t3 = _context2.v;
              if (_t3.name !== 'AbortError') _this4.$message.error('Share failed: ' + _t3.message);
            case 14:
              _context2.p = 14;
              loading.close();
              return _context2.f(14);
            case 15:
              return _context2.a(2);
          }
        }, _callee2, null, [[3, 9, 10, 11], [2, 13, 14, 15]]);
      }))();
    },
    editProduct: function editProduct(product) {
      this.$router.push({
        name: 'product-edit',
        params: {
          id: product.id
        }
      });
    },
    copyProductInfo: function copyProductInfo(product) {
      var _this5 = this;
      return _asyncToGenerator(/*#__PURE__*/_regenerator().m(function _callee3() {
        var sizes, text, _t4;
        return _regenerator().w(function (_context3) {
          while (1) switch (_context3.p = _context3.n) {
            case 0:
              sizes = _this5.formatSizes(product).join(' - ');
              text = "".concat(product.name, " ").concat(sizes);
              _context3.p = 1;
              _context3.n = 2;
              return navigator.clipboard.writeText(text);
            case 2:
              _this5.$message.success('Copied to clipboard!');
              _context3.n = 4;
              break;
            case 3:
              _context3.p = 3;
              _t4 = _context3.v;
              _this5.$message.error('Failed to copy');
            case 4:
              return _context3.a(2);
          }
        }, _callee3, null, [[1, 3]]);
      }))();
    }
  }, "editProduct", function editProduct(product) {
    this.editForm = {
      id: product.id,
      type: product.type,
      price: product.price,
      description: product.description,
      sizes: this.allSizes.map(function (name) {
        var existing = product.sizes.find(function (s) {
          return s.name === name;
        });
        return {
          name: name,
          value: existing ? existing.value : ''
        };
      }),
      image: null,
      imagePreview: product.path_thumb
    };
    this.editVisible = true;
  }), "handleEditFileChange", function handleEditFileChange(e) {
    var file = e.target.files[0];
    if (file) {
      this.editForm.image = file;
      this.editForm.imagePreview = URL.createObjectURL(file);
    }
  }), "parseEditQuickInput", function parseEditQuickInput(val) {
    var _this6 = this;
    if (!val) return;
    var nums = val.trim().split(/\s+/);
    var config = this.mappings[this.editForm.type];
    if (config && config.indices) {
      config.indices.forEach(function (idx, i) {
        if (nums[i]) _this6.editForm.sizes[idx].value = nums[i];
      });
    }
  }), "saveProduct", function saveProduct() {
    var _this7 = this;
    return _asyncToGenerator(/*#__PURE__*/_regenerator().m(function _callee4() {
      var formData, _t5;
      return _regenerator().w(function (_context4) {
        while (1) switch (_context4.p = _context4.n) {
          case 0:
            _this7.editLoading = true;
            _context4.p = 1;
            formData = new FormData();
            formData.append('type', _this7.editForm.type);
            formData.append('price', _this7.editForm.price);
            formData.append('description', _this7.editForm.description || '');
            if (_this7.editForm.image) {
              formData.append('banner', _this7.editForm.image);
            }
            _this7.editForm.sizes.forEach(function (s, i) {
              formData.append("sizes[".concat(i, "][name]"), s.name);
              formData.append("sizes[".concat(i, "][value]"), s.value);
            });
            _context4.n = 2;
            return axios.post("/salev2/api/product/".concat(_this7.editForm.id, "/update"), formData, {
              headers: {
                'Content-Type': 'multipart/form-data'
              }
            });
          case 2:
            _this7.$message.success('Product Updated');
            _this7.editVisible = false;
            _this7.fetchProducts(_this7.page);
            _context4.n = 4;
            break;
          case 3:
            _context4.p = 3;
            _t5 = _context4.v;
            _this7.$message.error('Update failed');
          case 4:
            _context4.p = 4;
            _this7.editLoading = false;
            return _context4.f(4);
          case 5:
            return _context4.a(2);
        }
      }, _callee4, null, [[1, 3, 4, 5]]);
    }))();
  })
});

/***/ }),

/***/ "./node_modules/babel-loader/lib/index.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/SmartUpload.vue?vue&type=script&lang=js":
/*!****************************************************************************************************************************************************************************!*\
  !*** ./node_modules/babel-loader/lib??ref--4-0!./node_modules/vue-loader/lib??vue-loader-options!./resources/js/components/salev2/SmartUpload.vue?vue&type=script&lang=js ***!
  \****************************************************************************************************************************************************************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
function _regenerator() { /*! regenerator-runtime -- Copyright (c) 2014-present, Facebook, Inc. -- license (MIT): https://github.com/babel/babel/blob/main/packages/babel-helpers/LICENSE */ var e, t, r = "function" == typeof Symbol ? Symbol : {}, n = r.iterator || "@@iterator", o = r.toStringTag || "@@toStringTag"; function i(r, n, o, i) { var c = n && n.prototype instanceof Generator ? n : Generator, u = Object.create(c.prototype); return _regeneratorDefine2(u, "_invoke", function (r, n, o) { var i, c, u, f = 0, p = o || [], y = !1, G = { p: 0, n: 0, v: e, a: d, f: d.bind(e, 4), d: function d(t, r) { return i = t, c = 0, u = e, G.n = r, a; } }; function d(r, n) { for (c = r, u = n, t = 0; !y && f && !o && t < p.length; t++) { var o, i = p[t], d = G.p, l = i[2]; r > 3 ? (o = l === n) && (u = i[(c = i[4]) ? 5 : (c = 3, 3)], i[4] = i[5] = e) : i[0] <= d && ((o = r < 2 && d < i[1]) ? (c = 0, G.v = n, G.n = i[1]) : d < l && (o = r < 3 || i[0] > n || n > l) && (i[4] = r, i[5] = n, G.n = l, c = 0)); } if (o || r > 1) return a; throw y = !0, n; } return function (o, p, l) { if (f > 1) throw TypeError("Generator is already running"); for (y && 1 === p && d(p, l), c = p, u = l; (t = c < 2 ? e : u) || !y;) { i || (c ? c < 3 ? (c > 1 && (G.n = -1), d(c, u)) : G.n = u : G.v = u); try { if (f = 2, i) { if (c || (o = "next"), t = i[o]) { if (!(t = t.call(i, u))) throw TypeError("iterator result is not an object"); if (!t.done) return t; u = t.value, c < 2 && (c = 0); } else 1 === c && (t = i["return"]) && t.call(i), c < 2 && (u = TypeError("The iterator does not provide a '" + o + "' method"), c = 1); i = e; } else if ((t = (y = G.n < 0) ? u : r.call(n, G)) !== a) break; } catch (t) { i = e, c = 1, u = t; } finally { f = 1; } } return { value: t, done: y }; }; }(r, o, i), !0), u; } var a = {}; function Generator() {} function GeneratorFunction() {} function GeneratorFunctionPrototype() {} t = Object.getPrototypeOf; var c = [][n] ? t(t([][n]())) : (_regeneratorDefine2(t = {}, n, function () { return this; }), t), u = GeneratorFunctionPrototype.prototype = Generator.prototype = Object.create(c); function f(e) { return Object.setPrototypeOf ? Object.setPrototypeOf(e, GeneratorFunctionPrototype) : (e.__proto__ = GeneratorFunctionPrototype, _regeneratorDefine2(e, o, "GeneratorFunction")), e.prototype = Object.create(u), e; } return GeneratorFunction.prototype = GeneratorFunctionPrototype, _regeneratorDefine2(u, "constructor", GeneratorFunctionPrototype), _regeneratorDefine2(GeneratorFunctionPrototype, "constructor", GeneratorFunction), GeneratorFunction.displayName = "GeneratorFunction", _regeneratorDefine2(GeneratorFunctionPrototype, o, "GeneratorFunction"), _regeneratorDefine2(u), _regeneratorDefine2(u, o, "Generator"), _regeneratorDefine2(u, n, function () { return this; }), _regeneratorDefine2(u, "toString", function () { return "[object Generator]"; }), (_regenerator = function _regenerator() { return { w: i, m: f }; })(); }
function _regeneratorDefine2(e, r, n, t) { var i = Object.defineProperty; try { i({}, "", {}); } catch (e) { i = 0; } _regeneratorDefine2 = function _regeneratorDefine(e, r, n, t) { function o(r, n) { _regeneratorDefine2(e, r, function (e) { return this._invoke(r, n, e); }); } r ? i ? i(e, r, { value: n, enumerable: !t, configurable: !t, writable: !t }) : e[r] = n : (o("next", 0), o("throw", 1), o("return", 2)); }, _regeneratorDefine2(e, r, n, t); }
function asyncGeneratorStep(n, t, e, r, o, a, c) { try { var i = n[a](c), u = i.value; } catch (n) { return void e(n); } i.done ? t(u) : Promise.resolve(u).then(r, o); }
function _asyncToGenerator(n) { return function () { var t = this, e = arguments; return new Promise(function (r, o) { var a = n.apply(t, e); function _next(n) { asyncGeneratorStep(a, r, o, _next, _throw, "next", n); } function _throw(n) { asyncGeneratorStep(a, r, o, _next, _throw, "throw", n); } _next(void 0); }); }; }
/* harmony default export */ __webpack_exports__["default"] = ({
  name: 'SmartUpload',
  data: function data() {
    return {
      productType: 'TROUSERS',
      productTypes: ['TROUSERS', 'JACKET', 'SHIRT', 'BLAZER', 'TIE', 'GILE', 'BELT', 'POLO SHIRT', 'HAT', 'SUIT'],
      selectedFiles: [],
      uploading: false,
      progress: 0,
      sessionId: null,
      totalCount: 0,
      processingCount: 0,
      successCount: 0,
      failedCount: 0,
      uploadLog: [],
      processingLog: [],
      statusCheckInterval: null
    };
  },
  destroyed: function destroyed() {
    if (this.statusCheckInterval) clearInterval(this.statusCheckInterval);
  },
  methods: {
    triggerFileSelect: function triggerFileSelect() {
      this.$refs.fileInput.click();
    },
    handleFileChange: function handleFileChange(e) {
      this.selectedFiles = Array.from(e.target.files);
      this.totalCount = this.selectedFiles.length;
    },
    generateSessionId: function generateSessionId() {
      return 'upl_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9);
    },
    startUpload: function startUpload() {
      var _this = this;
      return _asyncToGenerator(/*#__PURE__*/_regenerator().m(function _callee() {
        var i;
        return _regenerator().w(function (_context) {
          while (1) switch (_context.n) {
            case 0:
              _this.sessionId = _this.generateSessionId();
              _this.uploading = true;
              _this.progress = 0;
              _this.uploadLog = [];
              _this.processingLog = [];
              _this.successCount = 0;
              _this.failedCount = 0;
              _this.processingCount = 0;
              _this.startStatusCheck();
              i = 0;
            case 1:
              if (!(i < _this.selectedFiles.length)) {
                _context.n = 4;
                break;
              }
              _context.n = 2;
              return _this.uploadSingleFile(_this.selectedFiles[i], i + 1);
            case 2:
              _this.progress = Math.round((i + 1) / _this.selectedFiles.length * 100);
            case 3:
              i++;
              _context.n = 1;
              break;
            case 4:
              _this.uploading = false;
            case 5:
              return _context.a(2);
          }
        }, _callee);
      }))();
    },
    uploadSingleFile: function uploadSingleFile(file, index) {
      var _this2 = this;
      return _asyncToGenerator(/*#__PURE__*/_regenerator().m(function _callee2() {
        var formData, logItem, res, _err$response, _t;
        return _regenerator().w(function (_context2) {
          while (1) switch (_context2.p = _context2.n) {
            case 0:
              formData = new FormData();
              formData.append('image', file);
              formData.append('type', _this2.productType);
              formData.append('session_id', _this2.sessionId);
              logItem = {
                index: index,
                filename: file.name,
                status: 'uploading',
                message: 'Sending file...'
              };
              _this2.uploadLog.unshift(logItem);
              _context2.p = 1;
              _context2.n = 2;
              return axios.post('/sale/product-batch-upload/single', formData, {
                headers: {
                  'Content-Type': 'multipart/form-data'
                }
              });
            case 2:
              res = _context2.v;
              if (res.data.success) {
                logItem.status = 'success';
                logItem.message = 'Queued for processing';
              } else {
                logItem.status = 'failed';
                logItem.message = res.data.message || 'Upload failed';
              }
              _context2.n = 4;
              break;
            case 3:
              _context2.p = 3;
              _t = _context2.v;
              logItem.status = 'failed';
              logItem.message = ((_err$response = _t.response) === null || _err$response === void 0 || (_err$response = _err$response.data) === null || _err$response === void 0 ? void 0 : _err$response.message) || _t.message;
            case 4:
              return _context2.a(2);
          }
        }, _callee2, null, [[1, 3]]);
      }))();
    },
    startStatusCheck: function startStatusCheck() {
      var _this3 = this;
      if (this.statusCheckInterval) clearInterval(this.statusCheckInterval);
      this.statusCheckInterval = setInterval(/*#__PURE__*/_asyncToGenerator(/*#__PURE__*/_regenerator().m(function _callee3() {
        var res, _t2;
        return _regenerator().w(function (_context3) {
          while (1) switch (_context3.p = _context3.n) {
            case 0:
              _context3.p = 0;
              _context3.n = 1;
              return axios.get('/sale/product-batch-upload/status', {
                params: {
                  session_id: _this3.sessionId
                }
              });
            case 1:
              res = _context3.v;
              _this3.successCount = res.data.success;
              _this3.failedCount = res.data.failed;
              _this3.processingCount = res.data.processing;
              _this3.processingLog = res.data.logs;
              if (_this3.processingCount === 0 && res.data.total >= _this3.totalCount && !_this3.uploading) {
                clearInterval(_this3.statusCheckInterval);
              }
              _context3.n = 3;
              break;
            case 2:
              _context3.p = 2;
              _t2 = _context3.v;
              console.error('Status check failed:', _t2);
            case 3:
              return _context3.a(2);
          }
        }, _callee3, null, [[0, 2]]);
      })), 2000);
    },
    logRowClass: function logRowClass(_ref2) {
      var row = _ref2.row;
      return row.status === 'success' ? 'success-row' : row.status === 'failed' ? 'danger-row' : '';
    },
    getProcStatusType: function getProcStatusType(s) {
      return {
        'SUCCESS': 'success',
        'FAILED': 'danger',
        'ERROR': 'danger',
        'PENDING': 'warning'
      }[s] || 'info';
    }
  }
});

/***/ }),

/***/ "./node_modules/babel-loader/lib/index.js?!./node_modules/vue-loader/lib/loaders/templateLoader.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/CustomerManager.vue?vue&type=template&id=1e0b8b4e":
/*!******************************************************************************************************************************************************************************************************************************************************!*\
  !*** ./node_modules/babel-loader/lib??ref--4-0!./node_modules/vue-loader/lib/loaders/templateLoader.js??ref--6!./node_modules/vue-loader/lib??vue-loader-options!./resources/js/components/salev2/CustomerManager.vue?vue&type=template&id=1e0b8b4e ***!
  \******************************************************************************************************************************************************************************************************************************************************/
/*! exports provided: render, staticRenderFns */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "render", function() { return render; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "staticRenderFns", function() { return staticRenderFns; });
var render = function render() {
  var _vm = this,
    _c = _vm._self._c;
  return _c("div", {
    staticClass: "px-4 md:px-6 py-8"
  }, [_c("div", {
    staticClass: "flex justify-between items-center mb-6 mobile-hide"
  }, [_vm._m(0), _vm._v(" "), _c("div", {
    staticClass: "flex gap-3"
  }, [_c("el-button", {
    attrs: {
      type: "primary",
      icon: "el-icon-plus"
    }
  }, [_vm._v("Add Customer")])], 1)]), _vm._v(" "), _c("div", {
    staticClass: "md:hidden flex justify-between items-center mb-6"
  }, [_c("h1", {
    staticClass: "text-2xl font-bold text-slate-800"
  }, [_vm._v("Customers")]), _vm._v(" "), _c("el-button", {
    attrs: {
      size: "small",
      type: "primary",
      icon: "el-icon-plus",
      circle: ""
    }
  })], 1), _vm._v(" "), _c("div", {
    staticClass: "bg-white p-4 rounded-2xl shadow-sm border border-slate-100 mb-6"
  }, [_c("el-input", {
    staticClass: "w-full md:w-96",
    attrs: {
      placeholder: "Search by Name, Phone, Address...",
      "prefix-icon": "el-icon-search",
      clearable: ""
    },
    on: {
      input: _vm.debouncedFetch
    },
    model: {
      value: _vm.search,
      callback: function callback($$v) {
        _vm.search = $$v;
      },
      expression: "search"
    }
  })], 1), _vm._v(" "), _c("div", {
    directives: [{
      name: "loading",
      rawName: "v-loading",
      value: _vm.loading,
      expression: "loading"
    }]
  }, [_c("div", {
    staticClass: "grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4"
  }, _vm._l(_vm.consumers, function (customer) {
    return _c("div", {
      key: customer.id,
      staticClass: "bg-white rounded-xl p-5 border border-slate-100 shadow-sm hover:shadow-md transition-shadow relative group"
    }, [_c("div", {
      staticClass: "flex justify-between items-start mb-3"
    }, [_c("div", {
      staticClass: "flex items-center gap-3"
    }, [_c("div", {
      staticClass: "w-10 h-10 rounded-full bg-slate-100 flex items-center justify-center text-slate-500 font-bold"
    }, [_vm._v("\n                          " + _vm._s(customer.name.charAt(0).toUpperCase()) + "\n                      ")]), _vm._v(" "), _c("div", [_c("h4", {
      staticClass: "font-bold text-slate-800 text-sm"
    }, [_vm._v(_vm._s(customer.name))]), _vm._v(" "), _c("p", {
      staticClass: "text-xs text-slate-500"
    }, [_vm._v(_vm._s(customer.phone))])])]), _vm._v(" "), _c("div", {
      staticClass: "bg-blue-50 text-blue-600 px-2 py-1 rounded text-[10px] font-bold uppercase tracking-wider"
    }, [_vm._v("\n                      " + _vm._s(customer.orders_count) + " Orders\n                  ")])]), _vm._v(" "), _c("div", {
      staticClass: "mb-4"
    }, [_c("p", {
      staticClass: "text-xs text-slate-500 line-clamp-2 min-h-[2.5em] flex items-center"
    }, [_c("i", {
      staticClass: "el-icon-location-outline mr-1"
    }), _vm._v(" \n                      " + _vm._s(customer.address || "No address provided") + "\n                  ")])]), _vm._v(" "), _c("div", {
      staticClass: "flex gap-2"
    }, [_c("el-button", {
      staticClass: "flex-1",
      attrs: {
        type: "primary",
        plain: "",
        size: "small",
        icon: "el-icon-shopping-cart-full"
      },
      on: {
        click: function click($event) {
          return _vm.createOrder(customer);
        }
      }
    }, [_vm._v("Create Order")]), _vm._v(" "), _c("el-button", {
      attrs: {
        size: "small",
        icon: "el-icon-time",
        circle: ""
      },
      on: {
        click: function click($event) {
          return _vm.$message.info("History View Coming Soon");
        }
      }
    })], 1)]);
  }), 0), _vm._v(" "), !_vm.consumers.length && !_vm.loading ? _c("div", {
    staticClass: "text-center py-20 bg-white rounded-3xl border border-dashed border-slate-200"
  }, [_c("i", {
    staticClass: "el-icon-user text-4xl text-slate-200 mb-4 block"
  }), _vm._v(" "), _c("p", {
    staticClass: "text-slate-400 font-medium"
  }, [_vm._v("No customers found")])]) : _vm._e(), _vm._v(" "), _c("div", {
    staticClass: "mt-8 flex justify-center"
  }, [_c("el-pagination", {
    attrs: {
      background: "",
      layout: "prev, pager, next",
      "current-page": _vm.page,
      total: _vm.total,
      "page-size": 20
    },
    on: {
      "update:currentPage": function updateCurrentPage($event) {
        _vm.page = $event;
      },
      "update:current-page": function updateCurrentPage($event) {
        _vm.page = $event;
      },
      "current-change": _vm.fetchConsumers
    }
  })], 1)])]);
};
var staticRenderFns = [function () {
  var _vm = this,
    _c = _vm._self._c;
  return _c("div", [_c("h1", {
    staticClass: "text-3xl font-bold text-gray-900 tracking-tight"
  }, [_vm._v("Customer Network")]), _vm._v(" "), _c("p", {
    staticClass: "text-slate-500 mt-1"
  }, [_vm._v("Manage profiles and purchase history")])]);
}];
render._withStripped = true;


/***/ }),

/***/ "./node_modules/babel-loader/lib/index.js?!./node_modules/vue-loader/lib/loaders/templateLoader.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/Dashboard.vue?vue&type=template&id=16392a9a":
/*!************************************************************************************************************************************************************************************************************************************************!*\
  !*** ./node_modules/babel-loader/lib??ref--4-0!./node_modules/vue-loader/lib/loaders/templateLoader.js??ref--6!./node_modules/vue-loader/lib??vue-loader-options!./resources/js/components/salev2/Dashboard.vue?vue&type=template&id=16392a9a ***!
  \************************************************************************************************************************************************************************************************************************************************/
/*! exports provided: render, staticRenderFns */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "render", function() { return render; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "staticRenderFns", function() { return staticRenderFns; });
var render = function render() {
  var _vm$data, _vm$data4, _vm$data6, _vm$data8;
  var _vm = this,
    _c = _vm._self._c;
  return _c("div", {
    staticClass: "px-6 py-8"
  }, [_c("div", {
    staticClass: "flex justify-between items-center mb-8 mobile-hide"
  }, [_vm._m(0), _vm._v(" "), _c("el-button", {
    attrs: {
      type: "primary",
      icon: "el-icon-refresh",
      loading: _vm.loading
    },
    on: {
      click: _vm.fetchStats
    }
  }, [_vm._v("Refresh Data")])], 1), _vm._v(" "), _c("div", {
    staticClass: "grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8"
  }, _vm._l(_vm.statCards, function (card) {
    return _c("div", {
      key: card.title,
      staticClass: "relative overflow-hidden bg-white rounded-2xl p-6 shadow-sm border border-slate-100 hover:shadow-md transition-shadow"
    }, [_c("div", {
      "class": ["absolute top-0 right-0 w-24 h-24 -mr-8 -mt-8 rounded-full opacity-10", card.bgClass]
    }), _vm._v(" "), _c("div", {
      staticClass: "relative z-10 flex items-center justify-between"
    }, [_c("div", [_c("p", {
      staticClass: "text-[10px] font-black text-slate-400 uppercase tracking-widest mb-1"
    }, [_vm._v(_vm._s(card.title))]), _vm._v(" "), _c("h3", {
      staticClass: "text-2xl font-black text-slate-900"
    }, [_vm._v(_vm._s(card.value))])]), _vm._v(" "), _c("div", {
      "class": ["w-12 h-12 rounded-xl flex items-center justify-center text-xl text-white shadow-lg", card.bgClass]
    }, [_c("i", {
      "class": card.icon
    })])])]);
  }), 0), _vm._v(" "), _c("div", {
    staticClass: "grid grid-cols-1 xl:grid-cols-3 gap-8 mb-8"
  }, [_c("div", {
    staticClass: "bg-white rounded-2xl shadow-sm border border-slate-100 p-6"
  }, [_c("div", {
    staticClass: "flex justify-between items-center mb-6"
  }, [_c("h3", {
    staticClass: "text-lg font-black text-slate-900"
  }, [_vm._v("All Time")]), _vm._v(" "), _c("div", {
    staticClass: "text-right"
  }, [_c("p", {
    staticClass: "text-[10px] font-bold text-slate-400 uppercase leading-none"
  }, [_vm._v("Total Revenue")]), _vm._v(" "), _c("p", {
    staticClass: "text-blue-600 font-black"
  }, [_vm._v(_vm._s(_vm.formatPrice((_vm$data = _vm.data) === null || _vm$data === void 0 ? void 0 : _vm$data.all.total_revenue)))])])]), _vm._v(" "), _c("el-table", {
    staticClass: "legacy-dash-table",
    staticStyle: {
      width: "100%"
    },
    attrs: {
      data: _vm.allTypeStats,
      size: "mini"
    }
  }, [_c("el-table-column", {
    attrs: {
      prop: "type",
      label: "Type"
    }
  }), _vm._v(" "), _c("el-table-column", {
    attrs: {
      label: "Sold",
      width: "70",
      align: "center"
    },
    scopedSlots: _vm._u([{
      key: "default",
      fn: function fn(scope) {
        var _vm$data2;
        return [_vm._v(_vm._s(_vm.getCount((_vm$data2 = _vm.data) === null || _vm$data2 === void 0 ? void 0 : _vm$data2.all.type_stats, scope.row.type, "SOLD")))];
      }
    }])
  }), _vm._v(" "), _c("el-table-column", {
    attrs: {
      label: "Avail",
      width: "70",
      align: "center"
    },
    scopedSlots: _vm._u([{
      key: "default",
      fn: function fn(scope) {
        var _vm$data3;
        return [_vm._v(_vm._s(_vm.getCount((_vm$data3 = _vm.data) === null || _vm$data3 === void 0 ? void 0 : _vm$data3.all.type_stats, scope.row.type, "AVAILABLE")))];
      }
    }])
  })], 1)], 1), _vm._v(" "), _c("div", {
    staticClass: "bg-white rounded-2xl shadow-sm border border-slate-100 p-6"
  }, [_c("div", {
    staticClass: "flex justify-between items-center mb-6"
  }, [_c("h3", {
    staticClass: "text-lg font-black text-slate-900"
  }, [_vm._v("Current Month")]), _vm._v(" "), _c("div", {
    staticClass: "text-right"
  }, [_c("p", {
    staticClass: "text-[10px] font-bold text-slate-400 uppercase leading-none"
  }, [_vm._v("Revenue")]), _vm._v(" "), _c("p", {
    staticClass: "text-emerald-600 font-black"
  }, [_vm._v(_vm._s(_vm.formatPrice((_vm$data4 = _vm.data) === null || _vm$data4 === void 0 ? void 0 : _vm$data4.current_month.total_revenue)))])])]), _vm._v(" "), _c("el-table", {
    staticClass: "legacy-dash-table",
    staticStyle: {
      width: "100%"
    },
    attrs: {
      data: _vm.allTypeStats,
      size: "mini"
    }
  }, [_c("el-table-column", {
    attrs: {
      prop: "type",
      label: "Type"
    }
  }), _vm._v(" "), _c("el-table-column", {
    attrs: {
      label: "Sold",
      width: "80",
      align: "right"
    },
    scopedSlots: _vm._u([{
      key: "default",
      fn: function fn(scope) {
        var _vm$data5;
        return [_vm._v(_vm._s(_vm.getCount((_vm$data5 = _vm.data) === null || _vm$data5 === void 0 ? void 0 : _vm$data5.current_month.type_stats, scope.row.type)))];
      }
    }])
  })], 1)], 1), _vm._v(" "), _c("div", {
    staticClass: "bg-white rounded-2xl shadow-sm border border-slate-100 p-6"
  }, [_c("div", {
    staticClass: "flex justify-between items-center mb-6"
  }, [_c("h3", {
    staticClass: "text-lg font-black text-slate-900"
  }, [_vm._v("Previous Month")]), _vm._v(" "), _c("div", {
    staticClass: "text-right"
  }, [_c("p", {
    staticClass: "text-[10px] font-bold text-slate-400 uppercase leading-none"
  }, [_vm._v("Revenue")]), _vm._v(" "), _c("p", {
    staticClass: "text-slate-600 font-black"
  }, [_vm._v(_vm._s(_vm.formatPrice((_vm$data6 = _vm.data) === null || _vm$data6 === void 0 ? void 0 : _vm$data6.prev_month.total_revenue)))])])]), _vm._v(" "), _c("el-table", {
    staticClass: "legacy-dash-table",
    staticStyle: {
      width: "100%"
    },
    attrs: {
      data: _vm.allTypeStats,
      size: "mini"
    }
  }, [_c("el-table-column", {
    attrs: {
      prop: "type",
      label: "Type"
    }
  }), _vm._v(" "), _c("el-table-column", {
    attrs: {
      label: "Sold",
      width: "80",
      align: "right"
    },
    scopedSlots: _vm._u([{
      key: "default",
      fn: function fn(scope) {
        var _vm$data7;
        return [_vm._v(_vm._s(_vm.getCount((_vm$data7 = _vm.data) === null || _vm$data7 === void 0 ? void 0 : _vm$data7.prev_month.type_stats, scope.row.type)))];
      }
    }])
  })], 1)], 1)]), _vm._v(" "), _c("div", {
    staticClass: "bg-white rounded-2xl shadow-sm border border-slate-100 p-8"
  }, [_c("h3", {
    staticClass: "text-xl font-black text-slate-900 mb-6"
  }, [_vm._v("Monthly Performance")]), _vm._v(" "), _c("el-table", {
    staticStyle: {
      width: "100%"
    },
    attrs: {
      data: (_vm$data8 = _vm.data) === null || _vm$data8 === void 0 ? void 0 : _vm$data8.monthly_sales
    }
  }, [_c("el-table-column", {
    attrs: {
      prop: "month",
      label: "Month",
      width: "150"
    },
    scopedSlots: _vm._u([{
      key: "default",
      fn: function fn(scope) {
        return [_c("span", {
          staticClass: "font-bold text-slate-700"
        }, [_vm._v(_vm._s(scope.row.month))])];
      }
    }])
  }), _vm._v(" "), _c("el-table-column", {
    attrs: {
      prop: "number",
      label: "Orders Count",
      width: "150",
      align: "center"
    }
  }), _vm._v(" "), _c("el-table-column", {
    attrs: {
      prop: "final_amount",
      label: "Total Revenue"
    },
    scopedSlots: _vm._u([{
      key: "default",
      fn: function fn(scope) {
        return [_c("span", {
          staticClass: "text-blue-600 font-black"
        }, [_vm._v(_vm._s(_vm.formatPrice(scope.row.final_amount)))])];
      }
    }])
  }), _vm._v(" "), _c("el-table-column", {
    attrs: {
      label: "Trend",
      width: "200"
    },
    scopedSlots: _vm._u([{
      key: "default",
      fn: function fn(scope) {
        return [_c("div", {
          staticClass: "w-full bg-slate-100 h-2 rounded-full overflow-hidden"
        }, [_c("div", {
          staticClass: "bg-blue-500 h-full",
          style: {
            width: _vm.getTrendWidth(scope.row.final_amount) + "%"
          }
        })])];
      }
    }])
  })], 1)], 1)]);
};
var staticRenderFns = [function () {
  var _vm = this,
    _c = _vm._self._c;
  return _c("div", [_c("h1", {
    staticClass: "text-3xl font-bold text-gray-900 tracking-tight"
  }, [_vm._v("Sales Dashboard")]), _vm._v(" "), _c("p", {
    staticClass: "text-slate-500 mt-1"
  }, [_vm._v("Real-time performance and inventory metrics")])]);
}];
render._withStripped = true;


/***/ }),

/***/ "./node_modules/babel-loader/lib/index.js?!./node_modules/vue-loader/lib/loaders/templateLoader.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/OrderCreate.vue?vue&type=template&id=5edbfd6e":
/*!**************************************************************************************************************************************************************************************************************************************************!*\
  !*** ./node_modules/babel-loader/lib??ref--4-0!./node_modules/vue-loader/lib/loaders/templateLoader.js??ref--6!./node_modules/vue-loader/lib??vue-loader-options!./resources/js/components/salev2/OrderCreate.vue?vue&type=template&id=5edbfd6e ***!
  \**************************************************************************************************************************************************************************************************************************************************/
/*! exports provided: render, staticRenderFns */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "render", function() { return render; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "staticRenderFns", function() { return staticRenderFns; });
var render = function render() {
  var _vm = this,
    _c = _vm._self._c;
  return _c("div", {
    staticClass: "px-4 md:px-6 py-6 h-[calc(100vh-80px)] overflow-hidden flex flex-col"
  }, [_c("div", {
    staticClass: "flex justify-between items-center mb-4 flex-none"
  }, [_c("div", {
    staticClass: "flex items-center gap-3"
  }, [_c("el-button", {
    attrs: {
      icon: "el-icon-back",
      circle: ""
    },
    on: {
      click: function click($event) {
        return _vm.$router.push("/salev2/customers");
      }
    }
  }), _vm._v(" "), _c("h1", {
    staticClass: "text-2xl font-bold text-slate-800"
  }, [_vm._v("Create Order")])], 1), _vm._v(" "), _c("div", {
    staticClass: "flex gap-2"
  }, [_c("el-button", {
    staticClass: "px-6 shadow-lg shadow-green-100",
    attrs: {
      type: "success",
      icon: "el-icon-check",
      loading: _vm.submitting
    },
    on: {
      click: _vm.submitOrder
    }
  }, [_vm._v("Complete Order")])], 1)]), _vm._v(" "), _c("div", {
    staticClass: "flex flex-col lg:flex-row gap-6 h-full overflow-hidden"
  }, [_c("div", {
    staticClass: "w-full lg:w-96 flex-none flex flex-col gap-4 overflow-y-auto pb-20 lg:pb-0"
  }, [_c("div", {
    staticClass: "bg-white p-4 rounded-xl border border-slate-100 shadow-sm"
  }, [_vm._m(0), _vm._v(" "), _vm.customer ? _c("div", {
    staticClass: "space-y-3"
  }, [_c("el-input", {
    attrs: {
      placeholder: "Name",
      size: "small"
    },
    model: {
      value: _vm.customer.name,
      callback: function callback($$v) {
        _vm.$set(_vm.customer, "name", $$v);
      },
      expression: "customer.name"
    }
  }, [_c("template", {
    slot: "prepend"
  }, [_vm._v("Name")])], 2), _vm._v(" "), _c("el-input", {
    attrs: {
      placeholder: "Phone",
      size: "small"
    },
    model: {
      value: _vm.customer.phone,
      callback: function callback($$v) {
        _vm.$set(_vm.customer, "phone", $$v);
      },
      expression: "customer.phone"
    }
  }, [_c("template", {
    slot: "prepend"
  }, [_vm._v("Phone")])], 2), _vm._v(" "), _c("el-input", {
    attrs: {
      placeholder: "Address",
      type: "textarea",
      rows: 2,
      size: "small"
    },
    model: {
      value: _vm.customer.address,
      callback: function callback($$v) {
        _vm.$set(_vm.customer, "address", $$v);
      },
      expression: "customer.address"
    }
  }), _vm._v(" "), _c("div", [_c("label", {
    staticClass: "text-xs font-bold text-slate-500 uppercase"
  }, [_vm._v("Note")]), _vm._v(" "), _c("el-input", {
    attrs: {
      type: "textarea",
      rows: 2,
      placeholder: "Internal note..."
    },
    model: {
      value: _vm.order.note,
      callback: function callback($$v) {
        _vm.$set(_vm.order, "note", $$v);
      },
      expression: "order.note"
    }
  })], 1)], 1) : _c("div", {
    staticClass: "text-center py-4"
  }, [_vm.loadingCustomer ? _c("el-spinner") : _c("p", {
    staticClass: "text-slate-400 text-sm"
  }, [_vm._v("No customer selected")])], 1)]), _vm._v(" "), _c("div", {
    staticClass: "bg-white p-4 rounded-xl border border-slate-100 shadow-sm flex-1"
  }, [_vm._m(1), _vm._v(" "), _c("div", {
    staticClass: "space-y-4"
  }, [_c("div", [_c("label", {
    staticClass: "text-xs font-bold text-slate-500 uppercase"
  }, [_vm._v("Public Description (In Receipt)")]), _vm._v(" "), _c("el-input", {
    attrs: {
      type: "textarea",
      rows: 2,
      placeholder: "Public note..."
    },
    model: {
      value: _vm.order.desc,
      callback: function callback($$v) {
        _vm.$set(_vm.order, "desc", $$v);
      },
      expression: "order.desc"
    }
  })], 1), _vm._v(" "), _c("div", {
    staticClass: "pt-4 border-t border-slate-100 space-y-2"
  }, [_c("div", {
    staticClass: "flex justify-between items-center text-sm"
  }, [_c("span", {
    staticClass: "text-slate-500"
  }, [_vm._v("Subtotal (" + _vm._s(_vm.selectedProducts.length) + " items)")]), _vm._v(" "), _c("span", {
    staticClass: "font-bold"
  }, [_vm._v(_vm._s(_vm.formatPrice(_vm.subtotal)))])]), _vm._v(" "), _c("div", {
    staticClass: "flex items-center justify-between gap-2"
  }, [_c("span", {
    staticClass: "text-slate-500 text-sm"
  }, [_vm._v("Additional Fee")]), _vm._v(" "), _c("el-input-number", {
    staticClass: "w-28",
    attrs: {
      size: "mini",
      step: 10000
    },
    model: {
      value: _vm.order.additional_amount,
      callback: function callback($$v) {
        _vm.$set(_vm.order, "additional_amount", $$v);
      },
      expression: "order.additional_amount"
    }
  })], 1), _vm._v(" "), _c("div", {
    staticClass: "flex items-center justify-between gap-2"
  }, [_c("span", {
    staticClass: "text-slate-500 text-sm"
  }, [_vm._v("Discount / Sale")]), _vm._v(" "), _c("el-input-number", {
    staticClass: "w-28",
    attrs: {
      size: "mini",
      step: 10000
    },
    model: {
      value: _vm.order.freeship_amount,
      callback: function callback($$v) {
        _vm.$set(_vm.order, "freeship_amount", $$v);
      },
      expression: "order.freeship_amount"
    }
  })], 1), _vm._v(" "), _c("div", {
    staticClass: "flex items-center justify-between gap-2 pt-2 border-t border-dashed border-slate-100"
  }, [_c("span", {
    staticClass: "text-slate-700 text-sm font-medium"
  }, [_vm._v("Deposit")]), _vm._v(" "), _c("el-input-number", {
    staticClass: "w-28",
    attrs: {
      size: "mini",
      step: 50000
    },
    model: {
      value: _vm.order.deposit_amount,
      callback: function callback($$v) {
        _vm.$set(_vm.order, "deposit_amount", $$v);
      },
      expression: "order.deposit_amount"
    }
  })], 1), _vm._v(" "), _c("div", {
    staticClass: "flex justify-between items-center py-2"
  }, [_c("span", {
    staticClass: "font-bold text-slate-700 text-sm"
  }, [_vm._v("Paid in Full")]), _vm._v(" "), _c("el-switch", {
    model: {
      value: _vm.order.is_paid,
      callback: function callback($$v) {
        _vm.$set(_vm.order, "is_paid", $$v);
      },
      expression: "order.is_paid"
    }
  })], 1), _vm._v(" "), _c("div", {
    staticClass: "flex justify-between items-center pt-3 border-t border-slate-900"
  }, [_c("span", {
    staticClass: "font-black text-slate-900 text-lg"
  }, [_vm._v("TOTAL")]), _vm._v(" "), _c("span", {
    staticClass: "font-black text-blue-600 text-xl"
  }, [_vm._v(_vm._s(_vm.formatPrice(_vm.finalTotal)))])])])])])]), _vm._v(" "), _c("div", {
    staticClass: "flex-1 bg-white rounded-xl border border-slate-100 shadow-sm flex flex-col overflow-hidden"
  }, [_c("div", {
    staticClass: "p-4 border-b border-slate-100 flex gap-2"
  }, [_c("el-input", {
    staticClass: "flex-1",
    attrs: {
      placeholder: "Search product to add...",
      "prefix-icon": "el-icon-search",
      clearable: ""
    },
    on: {
      input: _vm.debouncedSearch
    },
    model: {
      value: _vm.productSearch,
      callback: function callback($$v) {
        _vm.productSearch = $$v;
      },
      expression: "productSearch"
    }
  }), _vm._v(" "), _c("el-button", {
    attrs: {
      icon: "el-icon-refresh",
      circle: ""
    },
    on: {
      click: _vm.searchProducts
    }
  })], 1), _vm._v(" "), _c("div", {
    staticClass: "flex-1 overflow-y-auto p-4 relative"
  }, [_vm.loadingProducts ? _c("div", {
    staticClass: "absolute inset-0 bg-white/80 z-10 flex items-center justify-center"
  }, [_c("el-spinner")], 1) : _vm._e(), _vm._v(" "), _vm.searchResults.length ? _c("div", {
    staticClass: "grid grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-3"
  }, _vm._l(_vm.searchResults, function (prod) {
    return _c("div", {
      key: prod.id,
      staticClass: "border rounded-lg p-2 cursor-pointer transition-all hover:shadow-md h-full flex flex-col",
      "class": _vm.isSelected(prod) ? "border-blue-500 bg-blue-50 ring-1 ring-blue-500" : "border-slate-100",
      on: {
        click: function click($event) {
          return _vm.toggleProduct(prod);
        }
      }
    }, [_c("div", {
      staticClass: "aspect-[3/4] bg-slate-100 rounded mb-2 overflow-hidden"
    }, [prod.path_thumb ? _c("img", {
      staticClass: "w-full h-full object-contain",
      attrs: {
        src: prod.path_thumb
      }
    }) : _vm._e()]), _vm._v(" "), _c("div", {
      staticClass: "text-center mt-auto"
    }, [_c("h5", {
      staticClass: "font-bold text-xs truncate"
    }, [_vm._v(_vm._s(prod.name))]), _vm._v(" "), _c("p", {
      staticClass: "text-xs text-blue-600 font-bold"
    }, [_vm._v(_vm._s(_vm.formatPrice(prod.price)))]), _vm._v(" "), _c("span", {
      staticClass: "text-[10px] text-slate-400 bg-slate-100 px-1 rounded"
    }, [_vm._v(_vm._s(prod.status))])])]);
  }), 0) : _vm.productSearch ? _c("div", {
    staticClass: "text-center py-10 text-slate-400"
  }, [_vm._v('\n                  No products found for "' + _vm._s(_vm.productSearch) + '"\n              ')]) : _c("div", {
    staticClass: "text-center py-10 text-slate-400"
  }, [_c("i", {
    staticClass: "el-icon-search text-4xl mb-2"
  }), _vm._v(" "), _c("p", [_vm._v("Type to search products")])])]), _vm._v(" "), _c("div", {
    staticClass: "p-3 bg-slate-50 border-t border-slate-200 h-24 overflow-x-auto whitespace-nowrap flex items-center gap-2"
  }, [_c("span", {
    staticClass: "text-xs font-bold text-slate-400 uppercase mr-2 flex-none"
  }, [_vm._v("Selected (" + _vm._s(_vm.selectedProducts.length) + ")")]), _vm._v(" "), !_vm.selectedProducts.length ? _c("div", {
    staticClass: "text-slate-400 text-sm italic"
  }, [_vm._v("No items selected")]) : _vm._e(), _vm._v(" "), _vm._l(_vm.selectedProducts, function (item) {
    return _c("div", {
      key: item.id,
      staticClass: "inline-flex items-center gap-2 bg-white border border-slate-200 rounded-lg p-1.5 shadow-sm pr-3"
    }, [item.path_thumb ? _c("img", {
      staticClass: "w-8 h-8 object-cover rounded",
      attrs: {
        src: item.path_thumb
      }
    }) : _vm._e(), _vm._v(" "), _c("div", {
      staticClass: "flex flex-col"
    }, [_c("span", {
      staticClass: "text-xs font-bold"
    }, [_vm._v(_vm._s(item.name))]), _vm._v(" "), _c("el-input-number", {
      staticClass: "w-20",
      attrs: {
        size: "mini",
        controls: false
      },
      model: {
        value: item.price,
        callback: function callback($$v) {
          _vm.$set(item, "price", $$v);
        },
        expression: "item.price"
      }
    })], 1), _vm._v(" "), _c("i", {
      staticClass: "el-icon-close text-slate-400 hover:text-red-500 cursor-pointer p-1",
      on: {
        click: function click($event) {
          $event.stopPropagation();
          return _vm.toggleProduct(item);
        }
      }
    })]);
  })], 2)])])]);
};
var staticRenderFns = [function () {
  var _vm = this,
    _c = _vm._self._c;
  return _c("h3", {
    staticClass: "font-bold text-slate-700 mb-3 flex items-center gap-2"
  }, [_c("i", {
    staticClass: "el-icon-user"
  }), _vm._v(" Customer")]);
}, function () {
  var _vm = this,
    _c = _vm._self._c;
  return _c("h3", {
    staticClass: "font-bold text-slate-700 mb-3 flex items-center gap-2"
  }, [_c("i", {
    staticClass: "el-icon-money"
  }), _vm._v(" Payment")]);
}];
render._withStripped = true;


/***/ }),

/***/ "./node_modules/babel-loader/lib/index.js?!./node_modules/vue-loader/lib/loaders/templateLoader.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/OrderEdit.vue?vue&type=template&id=2ac7fed2":
/*!************************************************************************************************************************************************************************************************************************************************!*\
  !*** ./node_modules/babel-loader/lib??ref--4-0!./node_modules/vue-loader/lib/loaders/templateLoader.js??ref--6!./node_modules/vue-loader/lib??vue-loader-options!./resources/js/components/salev2/OrderEdit.vue?vue&type=template&id=2ac7fed2 ***!
  \************************************************************************************************************************************************************************************************************************************************/
/*! exports provided: render, staticRenderFns */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "render", function() { return render; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "staticRenderFns", function() { return staticRenderFns; });
var render = function render() {
  var _vm = this,
    _c = _vm._self._c;
  return _c("div", {
    staticClass: "px-4 md:px-6 py-8"
  }, [_c("div", {
    staticClass: "flex justify-between items-center mb-6"
  }, [_vm._m(0), _vm._v(" "), _c("el-button", {
    attrs: {
      icon: "el-icon-back",
      circle: ""
    },
    on: {
      click: function click($event) {
        return _vm.$router.push("/salev2/orders");
      }
    }
  })], 1), _vm._v(" "), _c("div", {
    directives: [{
      name: "loading",
      rawName: "v-loading",
      value: _vm.loading,
      expression: "loading"
    }]
  }, [_vm.order ? _c("div", {
    staticClass: "grid grid-cols-1 lg:grid-cols-3 gap-6"
  }, [_c("div", {
    staticClass: "space-y-6"
  }, [_c("div", {
    staticClass: "bg-white p-6 rounded-2xl shadow-sm border border-slate-100"
  }, [_c("h2", {
    staticClass: "font-bold text-slate-800 text-lg mb-4"
  }, [_vm._v("Customer Info")]), _vm._v(" "), _c("div", {
    staticClass: "space-y-4"
  }, [_c("div", [_c("label", {
    staticClass: "block text-slate-400 font-bold text-xs uppercase mb-1"
  }, [_vm._v("Name")]), _vm._v(" "), _c("el-input", {
    model: {
      value: _vm.order.name,
      callback: function callback($$v) {
        _vm.$set(_vm.order, "name", $$v);
      },
      expression: "order.name"
    }
  })], 1), _vm._v(" "), _c("div", [_c("label", {
    staticClass: "block text-slate-400 font-bold text-xs uppercase mb-1"
  }, [_vm._v("Phone")]), _vm._v(" "), _c("el-input", {
    model: {
      value: _vm.order.phone,
      callback: function callback($$v) {
        _vm.$set(_vm.order, "phone", $$v);
      },
      expression: "order.phone"
    }
  })], 1), _vm._v(" "), _c("div", [_c("label", {
    staticClass: "block text-slate-400 font-bold text-xs uppercase mb-1"
  }, [_vm._v("Address")]), _vm._v(" "), _c("el-input", {
    attrs: {
      type: "textarea",
      rows: 2
    },
    model: {
      value: _vm.order.address,
      callback: function callback($$v) {
        _vm.$set(_vm.order, "address", $$v);
      },
      expression: "order.address"
    }
  })], 1), _vm._v(" "), _c("div", [_c("label", {
    staticClass: "block text-slate-400 font-bold text-xs uppercase mb-1"
  }, [_vm._v("Note")]), _vm._v(" "), _c("el-input", {
    attrs: {
      type: "textarea",
      rows: 2,
      placeholder: "Internal note..."
    },
    model: {
      value: _vm.order.note,
      callback: function callback($$v) {
        _vm.$set(_vm.order, "note", $$v);
      },
      expression: "order.note"
    }
  })], 1), _vm._v(" "), _c("div", [_c("label", {
    staticClass: "block text-slate-400 font-bold text-xs uppercase mb-1"
  }, [_vm._v("Public Description")]), _vm._v(" "), _c("el-input", {
    attrs: {
      type: "textarea",
      rows: 2,
      placeholder: "Public note..."
    },
    model: {
      value: _vm.order.desc,
      callback: function callback($$v) {
        _vm.$set(_vm.order, "desc", $$v);
      },
      expression: "order.desc"
    }
  })], 1)])]), _vm._v(" "), _c("div", {
    staticClass: "bg-white p-6 rounded-2xl shadow-sm border border-slate-100"
  }, [_c("h2", {
    staticClass: "font-bold text-slate-800 text-lg mb-4"
  }, [_vm._v("Order Settings")]), _vm._v(" "), _c("div", {
    staticClass: "space-y-4"
  }, [_c("div", [_c("label", {
    staticClass: "block text-slate-400 font-bold text-xs uppercase mb-1"
  }, [_vm._v("Status")]), _vm._v(" "), _c("el-select", {
    staticClass: "w-full",
    model: {
      value: _vm.order.status,
      callback: function callback($$v) {
        _vm.$set(_vm.order, "status", $$v);
      },
      expression: "order.status"
    }
  }, [_c("el-option", {
    attrs: {
      value: "NEW",
      label: "NEW"
    }
  }), _vm._v(" "), _c("el-option", {
    attrs: {
      value: "ON HOLD",
      label: "ON HOLD"
    }
  }), _vm._v(" "), _c("el-option", {
    attrs: {
      value: "A SHIP NOW",
      label: "A SHIP NOW"
    }
  }), _vm._v(" "), _c("el-option", {
    attrs: {
      value: "COMPLETED",
      label: "COMPLETED"
    }
  }), _vm._v(" "), _c("el-option", {
    attrs: {
      value: "CANCELED",
      label: "CANCELED"
    }
  })], 1)], 1), _vm._v(" "), _c("div", {
    staticClass: "flex items-center gap-2"
  }, [_c("el-switch", {
    attrs: {
      "active-color": "#13ce66"
    },
    model: {
      value: _vm.order.is_freeship,
      callback: function callback($$v) {
        _vm.$set(_vm.order, "is_freeship", $$v);
      },
      expression: "order.is_freeship"
    }
  }), _vm._v(" "), _c("span", {
    staticClass: "text-sm font-medium"
  }, [_vm._v("Free Shipping (20k)")])], 1), _vm._v(" "), _c("div", {
    staticClass: "flex items-center gap-2"
  }, [_c("el-switch", {
    attrs: {
      "active-color": "#13ce66"
    },
    model: {
      value: _vm.order.is_paid,
      callback: function callback($$v) {
        _vm.$set(_vm.order, "is_paid", $$v);
      },
      expression: "order.is_paid"
    }
  }), _vm._v(" "), _c("span", {
    staticClass: "text-sm font-medium"
  }, [_vm._v("Paid in Full")])], 1)])])]), _vm._v(" "), _c("div", {
    staticClass: "lg:col-span-2 space-y-6"
  }, [_c("div", {
    staticClass: "bg-white p-6 rounded-2xl shadow-sm border border-slate-100"
  }, [_c("div", {
    staticClass: "flex justify-between items-center mb-4"
  }, [_c("h2", {
    staticClass: "font-bold text-slate-800 text-lg"
  }, [_vm._v("Products")]), _vm._v(" "), _c("el-button", {
    attrs: {
      type: "text",
      icon: "el-icon-plus"
    },
    on: {
      click: function click($event) {
        _vm.showProductSelector = true;
      }
    }
  }, [_vm._v("Add Product")])], 1), _vm._v(" "), _c("div", {
    staticClass: "space-y-3"
  }, [_vm._l(_vm.order.products, function (p, index) {
    return _c("div", {
      key: p.id,
      staticClass: "flex gap-4 p-3 bg-slate-50 rounded-xl border border-slate-100 relative group"
    }, [_c("img", {
      staticClass: "w-16 h-20 object-contain bg-white rounded border border-slate-200",
      attrs: {
        src: p.path_thumb
      }
    }), _vm._v(" "), _c("div", {
      staticClass: "flex-1"
    }, [_c("div", {
      staticClass: "flex justify-between items-start"
    }, [_c("div", [_c("h4", {
      staticClass: "font-bold text-slate-800"
    }, [_vm._v(_vm._s(p.name))]), _vm._v(" "), _c("span", {
      staticClass: "text-[10px] text-slate-500 font-bold uppercase"
    }, [_vm._v(_vm._s(p.type))])]), _vm._v(" "), _c("el-input-number", {
      staticClass: "w-28",
      attrs: {
        step: 10000,
        size: "mini",
        "controls-position": "right"
      },
      model: {
        value: p.price,
        callback: function callback($$v) {
          _vm.$set(p, "price", $$v);
        },
        expression: "p.price"
      }
    })], 1), _vm._v(" "), _c("div", {
      staticClass: "mt-2 text-xs text-slate-500"
    }, _vm._l(p.sizes, function (s) {
      return _c("span", {
        key: s.name,
        staticClass: "mr-2"
      }, [_vm._v(_vm._s(s.name) + ": " + _vm._s(Math.floor(s.value)))]);
    }), 0)]), _vm._v(" "), _c("button", {
      staticClass: "absolute -top-2 -right-2 bg-red-500 text-white rounded-full w-6 h-6 flex items-center justify-center shadow-sm opacity-0 group-hover:opacity-100 transition-opacity",
      on: {
        click: function click($event) {
          return _vm.removeProduct(index);
        }
      }
    }, [_c("i", {
      staticClass: "el-icon-close font-bold text-xs"
    })])]);
  }), _vm._v(" "), !_vm.order.products.length ? _c("div", {
    staticClass: "text-center py-8 text-slate-400 italic"
  }, [_vm._v("No products selected")]) : _vm._e()], 2), _vm._v(" "), _c("div", {
    staticClass: "mt-6 pt-6 border-t border-slate-100 space-y-3"
  }, [_c("div", {
    staticClass: "flex justify-between text-sm"
  }, [_c("span", {
    staticClass: "text-slate-500"
  }, [_vm._v("Subtotal")]), _vm._v(" "), _c("span", {
    staticClass: "font-bold"
  }, [_vm._v(_vm._s(_vm.formatPrice(_vm.subtotal)))])]), _vm._v(" "), _c("div", {
    staticClass: "flex justify-between text-sm items-center"
  }, [_c("span", {
    staticClass: "text-slate-500"
  }, [_vm._v("Additional Fee")]), _vm._v(" "), _c("el-input-number", {
    staticClass: "w-32",
    attrs: {
      step: 10000,
      size: "mini",
      "controls-position": "right"
    },
    model: {
      value: _vm.order.additional_amount,
      callback: function callback($$v) {
        _vm.$set(_vm.order, "additional_amount", $$v);
      },
      expression: "order.additional_amount"
    }
  })], 1), _vm._v(" "), _c("div", {
    staticClass: "flex justify-between text-sm items-center"
  }, [_c("span", {
    staticClass: "text-slate-500"
  }, [_vm._v("Deposit")]), _vm._v(" "), _c("el-input-number", {
    staticClass: "w-32",
    attrs: {
      step: 10000,
      size: "mini",
      "controls-position": "right"
    },
    model: {
      value: _vm.order.deposit_amount,
      callback: function callback($$v) {
        _vm.$set(_vm.order, "deposit_amount", $$v);
      },
      expression: "order.deposit_amount"
    }
  })], 1), _vm._v(" "), _c("div", {
    staticClass: "flex justify-between text-sm"
  }, [_c("span", {
    staticClass: "text-slate-500"
  }, [_vm._v("Shipping Fee")]), _vm._v(" "), _c("span", {
    staticClass: "font-bold"
  }, [_vm._v(_vm._s(_vm.order.is_freeship ? "0đ" : "20.000đ"))])]), _vm._v(" "), _c("div", {
    staticClass: "flex justify-between text-lg font-black text-slate-900 pt-3 border-t border-dashed border-slate-200"
  }, [_c("span", [_vm._v("Total")]), _vm._v(" "), _c("span", [_vm._v(_vm._s(_vm.formatPrice(_vm.totalAmount)))])]), _vm._v(" "), _c("div", {
    staticClass: "flex justify-between text-base font-bold text-blue-600"
  }, [_c("span", [_vm._v("Viettel Amount")]), _vm._v(" "), _c("span", [_vm._v(_vm._s(_vm.formatPrice(_vm.viettelAmount)))])])])]), _vm._v(" "), _c("div", {
    staticClass: "flex justify-end gap-3 pt-4"
  }, [_c("el-button", {
    on: {
      click: function click($event) {
        return _vm.$router.push("/salev2/orders");
      }
    }
  }, [_vm._v("Cancel")]), _vm._v(" "), _c("el-button", {
    staticClass: "px-8 shadow-lg shadow-blue-200",
    attrs: {
      type: "primary",
      loading: _vm.submitting
    },
    on: {
      click: _vm.saveOrder
    }
  }, [_vm._v("Save Changes")])], 1)])]) : _vm._e()]), _vm._v(" "), _c("el-dialog", {
    attrs: {
      title: "Add Products",
      visible: _vm.showProductSelector,
      width: "90%",
      "custom-class": "mobile-dialog"
    },
    on: {
      "update:visible": function updateVisible($event) {
        _vm.showProductSelector = $event;
      }
    }
  }, [_c("div", {
    staticClass: "mb-4"
  }, [_c("el-input", {
    attrs: {
      placeholder: "Search by name...",
      "prefix-icon": "el-icon-search"
    },
    on: {
      input: _vm.debouncedSearch
    },
    model: {
      value: _vm.productSearch,
      callback: function callback($$v) {
        _vm.productSearch = $$v;
      },
      expression: "productSearch"
    }
  })], 1), _vm._v(" "), _c("div", {
    staticClass: "grid grid-cols-2 md:grid-cols-4 gap-3 max-h-[60vh] overflow-y-auto"
  }, _vm._l(_vm.availableProducts, function (p) {
    return _c("div", {
      key: p.id,
      staticClass: "border rounded-lg p-2 cursor-pointer hover:border-blue-500 transition-colors relative",
      on: {
        click: function click($event) {
          return _vm.toggleProductSelect(p);
        }
      }
    }, [_c("img", {
      staticClass: "w-full aspect-[3/4] object-contain bg-slate-50 mb-2",
      attrs: {
        src: p.path_thumb
      }
    }), _vm._v(" "), _c("p", {
      staticClass: "font-bold text-xs text-center"
    }, [_vm._v(_vm._s(p.name))]), _vm._v(" "), _vm.isProductSelected(p.id) ? _c("div", {
      staticClass: "absolute inset-0 bg-blue-500/20 border-2 border-blue-500 rounded-lg flex items-center justify-center"
    }, [_c("i", {
      staticClass: "el-icon-check text-white bg-blue-500 rounded-full p-1"
    })]) : _vm._e()]);
  }), 0), _vm._v(" "), _c("span", {
    attrs: {
      slot: "footer"
    },
    slot: "footer"
  }, [_c("el-button", {
    on: {
      click: function click($event) {
        _vm.showProductSelector = false;
      }
    }
  }, [_vm._v("Done")])], 1)])], 1);
};
var staticRenderFns = [function () {
  var _vm = this,
    _c = _vm._self._c;
  return _c("div", [_c("h1", {
    staticClass: "text-3xl font-bold text-gray-900 tracking-tight"
  }, [_vm._v("Edit Order")]), _vm._v(" "), _c("p", {
    staticClass: "text-slate-500 mt-1"
  }, [_vm._v("Update order details, products, and status")])]);
}];
render._withStripped = true;


/***/ }),

/***/ "./node_modules/babel-loader/lib/index.js?!./node_modules/vue-loader/lib/loaders/templateLoader.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/OrderManager.vue?vue&type=template&id=64110eb0":
/*!***************************************************************************************************************************************************************************************************************************************************!*\
  !*** ./node_modules/babel-loader/lib??ref--4-0!./node_modules/vue-loader/lib/loaders/templateLoader.js??ref--6!./node_modules/vue-loader/lib??vue-loader-options!./resources/js/components/salev2/OrderManager.vue?vue&type=template&id=64110eb0 ***!
  \***************************************************************************************************************************************************************************************************************************************************/
/*! exports provided: render, staticRenderFns */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "render", function() { return render; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "staticRenderFns", function() { return staticRenderFns; });
var render = function render() {
  var _vm$currentCapture, _vm$currentCapture2, _vm$currentCapture3, _vm$currentCapture4, _vm$currentCapture5, _vm$currentCapture6, _vm$currentCapture7;
  var _vm = this,
    _c = _vm._self._c;
  return _c("div", {
    staticClass: "px-4 md:px-6 py-8"
  }, [_c("div", {
    staticClass: "flex flex-col md:flex-row justify-between items-center mb-6 gap-4"
  }, [_vm._m(0), _vm._v(" "), _c("div", {
    staticClass: "flex gap-2 w-full md:w-auto"
  }, [_c("el-input", {
    staticClass: "w-full md:w-64",
    attrs: {
      placeholder: "Phone, Name...",
      "prefix-icon": "el-icon-search"
    },
    on: {
      input: _vm.debouncedFetch
    },
    model: {
      value: _vm.search,
      callback: function callback($$v) {
        _vm.search = $$v;
      },
      expression: "search"
    }
  }), _vm._v(" "), _c("el-button", {
    attrs: {
      type: "primary",
      icon: "el-icon-plus"
    },
    on: {
      click: function click($event) {
        return _vm.$router.push("/salev2/order/create");
      }
    }
  }, [_vm._v("New")]), _vm._v(" "), _c("el-button", {
    attrs: {
      icon: "el-icon-refresh"
    },
    on: {
      click: function click($event) {
        return _vm.fetchOrders(1);
      }
    }
  })], 1)]), _vm._v(" "), _c("div", {
    staticClass: "mb-6 flex justify-between items-center"
  }, [_c("el-tabs", {
    staticClass: "premium-tabs",
    attrs: {
      type: "card"
    },
    on: {
      "tab-click": _vm.handleTabClick
    },
    model: {
      value: _vm.activeTab,
      callback: function callback($$v) {
        _vm.activeTab = $$v;
      },
      expression: "activeTab"
    }
  }, [_c("el-tab-pane", {
    attrs: {
      label: "New Orders",
      name: "NEW"
    }
  }), _vm._v(" "), _c("el-tab-pane", {
    attrs: {
      label: "All Orders",
      name: "ALL"
    }
  })], 1), _vm._v(" "), _c("div", {
    staticClass: "flex gap-2"
  }, [_c("el-radio-group", {
    attrs: {
      size: "small"
    },
    model: {
      value: _vm.viewMode,
      callback: function callback($$v) {
        _vm.viewMode = $$v;
      },
      expression: "viewMode"
    }
  }, [_c("el-radio-button", {
    attrs: {
      label: "grid"
    }
  }, [_c("i", {
    staticClass: "el-icon-menu"
  })]), _vm._v(" "), _c("el-radio-button", {
    attrs: {
      label: "table"
    }
  }, [_c("i", {
    staticClass: "el-icon-tickets"
  })])], 1)], 1)], 1), _vm._v(" "), _vm.viewMode === "grid" ? _c("div", {
    directives: [{
      name: "loading",
      rawName: "v-loading",
      value: _vm.loading,
      expression: "loading"
    }]
  }, [_c("div", {
    staticClass: "grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 xl:grid-cols-6 gap-4"
  }, _vm._l(_vm.orders, function (order) {
    return _c("div", {
      key: order.id,
      staticClass: "bg-white rounded-2xl border border-slate-100 shadow-sm hover:shadow-xl transition-all duration-300 flex flex-col overflow-hidden w-full md:max-w-[280px] mx-auto h-full"
    }, [_c("div", {
      staticClass: "p-4 border-b border-slate-50 bg-slate-50/30"
    }, [_c("div", {
      staticClass: "flex justify-between items-start mb-2"
    }, [_c("div", [_c("h5", {
      staticClass: "font-black text-slate-900 leading-none mb-1 uppercase tracking-tight"
    }, [_vm._v(_vm._s(order.name))]), _vm._v(" "), _c("div", {
      staticClass: "flex items-center gap-2"
    }, [_c("p", {
      staticClass: "text-blue-600 text-xs font-bold"
    }, [_vm._v(_vm._s(order.phone))]), _vm._v(" "), _c("span", {
      staticClass: "text-[10px] text-slate-400 font-medium italic"
    }, [_vm._v("• " + _vm._s(order.created_at_human))])])]), _vm._v(" "), _c("div", {
      staticClass: "flex items-center gap-1"
    }, [_c("el-tag", {
      staticClass: "rounded-pill uppercase font-bold text-[8px]",
      attrs: {
        type: _vm.getStatusType(order.status),
        size: "mini",
        effect: "dark"
      }
    }, [_vm._v("\n                      " + _vm._s(order.status) + "\n                  ")]), _vm._v(" "), _c("el-dropdown", {
      attrs: {
        trigger: "click"
      }
    }, [_c("span", {
      staticClass: "cursor-pointer p-1 rounded hover:bg-slate-100 transition-colors"
    }, [_c("i", {
      staticClass: "el-icon-more text-slate-400"
    })]), _vm._v(" "), _c("el-dropdown-menu", {
      attrs: {
        slot: "dropdown"
      },
      slot: "dropdown"
    }, [_c("el-dropdown-item", {
      attrs: {
        icon: "el-icon-camera"
      },
      nativeOn: {
        click: function click($event) {
          return _vm.captureOrder(order);
        }
      }
    }, [_vm._v("Capture Receipt")]), _vm._v(" "), _c("el-dropdown-item", {
      attrs: {
        icon: "el-icon-edit"
      },
      nativeOn: {
        click: function click($event) {
          return _vm.$router.push({
            name: "order-edit",
            params: {
              id: order.id
            }
          });
        }
      }
    }, [_vm._v("Edit Order")]), _vm._v(" "), _c("el-dropdown-item", {
      attrs: {
        icon: "el-icon-user"
      },
      nativeOn: {
        click: function click($event) {
          return _vm.viewCustomerHistory(order);
        }
      }
    }, [_vm._v("Customer History")]), _vm._v(" "), _c("el-divider", {
      staticClass: "my-1"
    }), _vm._v(" "), order.status !== "A SHIP NOW" ? _c("el-dropdown-item", {
      staticClass: "text-primary font-bold",
      attrs: {
        icon: "el-icon-truck"
      },
      nativeOn: {
        click: function click($event) {
          return _vm.updateStatus(order, "A SHIP NOW");
        }
      }
    }, [_vm._v("Ready to Ship")]) : _vm._e(), _vm._v(" "), order.status !== "ON HOLD" ? _c("el-dropdown-item", {
      staticClass: "text-warning font-bold",
      attrs: {
        icon: "el-icon-video-pause"
      },
      nativeOn: {
        click: function click($event) {
          return _vm.updateStatus(order, "ON HOLD");
        }
      }
    }, [_vm._v("Mark On Hold")]) : _vm._e(), _vm._v(" "), order.status !== "COMPLETED" ? _c("el-dropdown-item", {
      staticClass: "text-success font-bold",
      attrs: {
        icon: "el-icon-check"
      },
      nativeOn: {
        click: function click($event) {
          return _vm.updateStatus(order, "COMPLETED");
        }
      }
    }, [_vm._v("Complete")]) : _vm._e(), _vm._v(" "), _c("el-dropdown-item", {
      staticClass: "text-danger",
      attrs: {
        icon: "el-icon-delete"
      },
      nativeOn: {
        click: function click($event) {
          return _vm.confirmDelete(order.id);
        }
      }
    }, [_vm._v("Delete")])], 1)], 1)], 1)]), _vm._v(" "), _c("p", {
      staticClass: "text-[12px] text-slate-500 break-words leading-snug"
    }, [_c("i", {
      staticClass: "el-icon-location-outline mr-1"
    }), _vm._v(_vm._s(order.address))]), _vm._v(" "), order.note ? _c("div", {
      staticClass: "mt-2 text-[10px] bg-red-50 text-red-600 p-2 rounded-lg font-bold"
    }, [_c("i", {
      staticClass: "el-icon-warning mr-1"
    }), _vm._v(_vm._s(order.note) + "\n          ")]) : _vm._e()]), _vm._v(" "), _c("div", {
      staticClass: "p-4 flex-1 space-y-4"
    }, [_c("div", {
      staticClass: "grid grid-cols-2 gap-2 pb-2"
    }, _vm._l(order.products, function (p) {
      return _c("img", {
        key: p.id,
        staticClass: "w-full h-auto object-contain bg-slate-100 rounded shadow-sm border border-slate-100",
        attrs: {
          src: p.path_thumb
        }
      });
    }), 0), _vm._v(" "), _c("div", {
      staticClass: "bg-gray-50 rounded-xl p-3 text-[10px] space-y-1"
    }, [_vm._l(order.products, function (p) {
      return _c("div", {
        key: p.id,
        staticClass: "text-slate-600 block mb-1"
      }, [_c("div", {
        staticClass: "font-bold flex justify-between"
      }, [_c("span", [_vm._v(_vm._s(p.type) + " " + _vm._s(p.name))]), _vm._v(" "), _c("span", {
        staticClass: "text-slate-400"
      }, [_vm._v(_vm._s(_vm.formatPrice(p.price)))])]), _vm._v(" "), _c("div", {
        staticClass: "flex flex-wrap gap-1 mt-0.5"
      }, _vm._l(_vm.formatSizes(p), function (size) {
        return _c("span", {
          key: size,
          staticClass: "text-[8px] bg-slate-200 text-slate-600 px-1 py-0.5 rounded leading-none uppercase"
        }, [_vm._v("\n                       " + _vm._s(size) + "\n                     ")]);
      }), 0)]);
    }), _vm._v(" "), order.additional_amount > 0 ? _c("div", {
      staticClass: "text-blue-600 font-bold border-t border-slate-100 pt-2 flex justify-between"
    }, [_c("span", [_vm._v("Phí cộng thêm:")]), _vm._v(" "), _c("span", [_vm._v("+" + _vm._s(_vm.formatPrice(order.additional_amount)))])]) : _vm._e(), _vm._v(" "), order.is_paid || order.deposit_amount > 0 ? _c("div", {
      staticClass: "text-emerald-600 font-bold flex justify-between"
    }, [_c("span", [_vm._v(_vm._s(order.is_paid ? "Đã thanh toán:" : "Cọc:"))]), _vm._v(" "), _c("span", [_vm._v("-" + _vm._s(_vm.formatPrice(order.is_paid ? order.total_amount : order.deposit_amount)))])]) : _vm._e()], 2)]), _vm._v(" "), _c("div", {
      staticClass: "p-4 pt-0"
    }, [_c("div", {
      staticClass: "flex justify-between items-end"
    }, [_c("div", [_c("span", {
      staticClass: "text-[9px] text-slate-400 font-black uppercase block leading-none mb-1"
    }, [_vm._v("Viettel Amount")]), _vm._v(" "), _c("div", {
      staticClass: "flex items-baseline gap-1"
    }, [_c("span", {
      staticClass: "text-xl font-black text-slate-900"
    }, [_vm._v(_vm._s(_vm.formatPrice(order.viettel_amount)))]), _vm._v(" "), _c("span", {
      staticClass: "text-[10px] font-bold text-slate-400"
    }, [_vm._v("+ phí ship")])])]), _vm._v(" "), order.is_old ? _c("div", {
      staticClass: "animate-pulse bg-red-500 text-white text-[8px] font-black px-2 py-0.5 rounded-pill"
    }, [_vm._v("DELAY")]) : _vm._e()])])]);
  }), 0)]) : _c("div", {
    directives: [{
      name: "loading",
      rawName: "v-loading",
      value: _vm.loading,
      expression: "loading"
    }],
    staticClass: "bg-white rounded-2xl shadow-sm border border-slate-100 overflow-hidden"
  }, [_c("el-table", {
    staticStyle: {
      width: "100%"
    },
    attrs: {
      data: _vm.orders
    }
  }, [_c("el-table-column", {
    attrs: {
      label: "Consumer",
      width: "200"
    },
    scopedSlots: _vm._u([{
      key: "default",
      fn: function fn(scope) {
        return [_c("div", {
          staticClass: "flex flex-col text-xs"
        }, [_c("span", {
          staticClass: "font-bold text-slate-800"
        }, [_vm._v(_vm._s(scope.row.name))]), _vm._v(" "), _c("span", {
          staticClass: "text-slate-500"
        }, [_vm._v(_vm._s(scope.row.phone))]), _vm._v(" "), _c("span", {
          staticClass: "text-[10px] text-slate-400 italic"
        }, [_vm._v(_vm._s(scope.row.created_at_human))])])];
      }
    }])
  }), _vm._v(" "), _c("el-table-column", {
    attrs: {
      label: "Summary",
      "min-width": "250"
    },
    scopedSlots: _vm._u([{
      key: "default",
      fn: function fn(scope) {
        return [_c("div", {
          staticClass: "text-[10px] space-y-0.5 text-slate-500 italic"
        }, [_vm._l(scope.row.products, function (p) {
          return _c("span", {
            key: p.id,
            staticClass: "block"
          }, [_vm._v("\n                  " + _vm._s(p.type) + " " + _vm._s(p.name) + ": " + _vm._s(_vm.formatPrice(p.price)) + "\n              ")]);
        }), _vm._v(" "), _c("span", {
          staticClass: "font-bold text-blue-600 block pt-1 border-t border-slate-50"
        }, [_vm._v("\n                  Total: " + _vm._s(_vm.formatPrice(scope.row.is_paid ? 0 : scope.row.total_amount - scope.row.deposit_amount)) + " \n                  " + _vm._s(scope.row.is_paid ? "(PAID)" : "+ phí ship") + "\n              ")])], 2)];
      }
    }])
  }), _vm._v(" "), _c("el-table-column", {
    attrs: {
      label: "Additionals",
      width: "120",
      align: "center"
    },
    scopedSlots: _vm._u([{
      key: "default",
      fn: function fn(scope) {
        return [_c("p", {
          staticClass: "text-[10px]"
        }, [_c("span", {
          staticClass: "text-blue-600"
        }, [_vm._v("Phí thêm:")]), _vm._v(" " + _vm._s(_vm.formatPrice(scope.row.additional_amount)))]), _vm._v(" "), _c("p", {
          staticClass: "text-[10px]"
        }, [_c("span", {
          staticClass: "text-emerald-600"
        }, [_vm._v("Cọc:")]), _vm._v(" " + _vm._s(_vm.formatPrice(scope.row.deposit_amount)))])];
      }
    }])
  }), _vm._v(" "), _c("el-table-column", {
    attrs: {
      label: "Viettel Amt",
      width: "120",
      align: "center"
    },
    scopedSlots: _vm._u([{
      key: "default",
      fn: function fn(scope) {
        return [_c("span", {
          staticClass: "bg-gray-800 text-white px-1.5 py-0.5 rounded text-[10px] font-bold"
        }, [_vm._v(_vm._s(_vm.formatPrice(scope.row.viettel_amount)))])];
      }
    }])
  }), _vm._v(" "), _c("el-table-column", {
    attrs: {
      label: "Status",
      width: "120"
    },
    scopedSlots: _vm._u([{
      key: "default",
      fn: function fn(scope) {
        return [_c("el-tag", {
          staticClass: "rounded-pill uppercase font-bold text-[10px]",
          attrs: {
            type: _vm.getStatusType(scope.row.status),
            size: "mini",
            effect: "dark"
          }
        }, [_vm._v("\n            " + _vm._s(scope.row.status) + "\n          ")])];
      }
    }])
  }), _vm._v(" "), _c("el-table-column", {
    attrs: {
      label: "Action",
      align: "right",
      width: "100"
    },
    scopedSlots: _vm._u([{
      key: "default",
      fn: function fn(scope) {
        return [_c("el-dropdown", {
          attrs: {
            trigger: "click"
          }
        }, [_c("el-button", {
          attrs: {
            type: "text",
            icon: "el-icon-more"
          }
        }), _vm._v(" "), _c("el-dropdown-menu", {
          attrs: {
            slot: "dropdown"
          },
          slot: "dropdown"
        }, [_c("el-dropdown-item", {
          attrs: {
            icon: "el-icon-camera"
          },
          nativeOn: {
            click: function click($event) {
              return _vm.captureOrder(scope.row);
            }
          }
        }, [_vm._v("Capture Order")]), _vm._v(" "), _c("el-dropdown-item", {
          attrs: {
            icon: "el-icon-view"
          },
          nativeOn: {
            click: function click($event) {
              return _vm.viewOrder(scope.row);
            }
          }
        }, [_vm._v("View Detail")]), _vm._v(" "), _c("el-dropdown-item", {
          staticClass: "text-danger",
          attrs: {
            icon: "el-icon-delete"
          },
          nativeOn: {
            click: function click($event) {
              return _vm.confirmDelete(scope.row.id);
            }
          }
        }, [_vm._v("Delete")])], 1)], 1)];
      }
    }])
  })], 1)], 1), _vm._v(" "), _vm.orders.length ? _c("div", {
    staticClass: "mt-8 flex justify-center"
  }, [_c("el-pagination", {
    attrs: {
      background: "",
      layout: "prev, pager, next",
      "current-page": _vm.page,
      total: _vm.total,
      "page-size": 10
    },
    on: {
      "update:currentPage": function updateCurrentPage($event) {
        _vm.page = $event;
      },
      "update:current-page": function updateCurrentPage($event) {
        _vm.page = $event;
      },
      "current-change": _vm.handlePageChange
    }
  })], 1) : _vm._e(), _vm._v(" "), _c("div", {
    staticStyle: {
      position: "absolute",
      left: "-9999px",
      top: "0"
    }
  }, [_c("div", {
    ref: "receiptTemplate",
    staticClass: "p-8 bg-white w-[450px] font-sans text-slate-800",
    attrs: {
      id: "order-receipt-v2"
    }
  }, [_vm._m(1), _vm._v(" "), _c("div", {
    staticClass: "bg-slate-50 rounded-2xl p-6 mb-5"
  }, [_c("div", {
    staticClass: "flex items-center mb-3"
  }, [_c("div", {
    staticClass: "bg-blue-600 text-white w-10 h-10 rounded-full flex items-center justify-center font-bold mr-3"
  }, [_vm._v("P")]), _vm._v(" "), _c("div", [_c("span", {
    staticClass: "text-[10px] text-slate-400 uppercase font-bold tracking-wider"
  }, [_vm._v("KHÁCH HÀNG")]), _vm._v(" "), _c("h2", {
    staticClass: "text-lg font-bold m-0"
  }, [_vm._v(_vm._s((_vm$currentCapture = _vm.currentCapture) === null || _vm$currentCapture === void 0 ? void 0 : _vm$currentCapture.name))])])]), _vm._v(" "), _c("div", {
    staticClass: "text-sm text-slate-600 space-y-1"
  }, [_c("p", [_c("i", {
    staticClass: "el-icon-phone mr-2"
  }), _vm._v(_vm._s((_vm$currentCapture2 = _vm.currentCapture) === null || _vm$currentCapture2 === void 0 ? void 0 : _vm$currentCapture2.phone))]), _vm._v(" "), _c("p", [_c("i", {
    staticClass: "el-icon-location mr-2"
  }), _vm._v(_vm._s((_vm$currentCapture3 = _vm.currentCapture) === null || _vm$currentCapture3 === void 0 ? void 0 : _vm$currentCapture3.address))])])]), _vm._v(" "), (_vm$currentCapture4 = _vm.currentCapture) !== null && _vm$currentCapture4 !== void 0 && _vm$currentCapture4.note ? _c("div", {
    staticClass: "bg-yellow-50 border border-yellow-100 rounded-xl p-4 mb-6"
  }, [_c("p", {
    staticClass: "text-sm font-medium"
  }, [_vm._v(_vm._s(_vm.currentCapture.note))])]) : _vm._e(), _vm._v(" "), _c("h3", {
    staticClass: "text-[10px] font-black text-slate-300 uppercase tracking-widest border-b border-slate-100 pb-2 mb-4"
  }, [_vm._v("ĐƠN HÀNG CỦA BẠN")]), _vm._v(" "), _c("div", {
    staticClass: "flex flex-wrap gap-4 mb-6"
  }, _vm._l((_vm$currentCapture5 = _vm.currentCapture) === null || _vm$currentCapture5 === void 0 ? void 0 : _vm$currentCapture5.products, function (p) {
    return _c("div", {
      key: p.id,
      staticClass: "w-[calc(50%-8px)] text-center"
    }, [_c("img", {
      staticClass: "w-full aspect-[3/4] object-cover rounded-xl bg-slate-50 mb-2 border border-slate-100",
      attrs: {
        src: p.path_thumb
      }
    }), _vm._v(" "), _c("span", {
      staticClass: "text-sm font-bold"
    }, [_vm._v(_vm._s(p.name))])]);
  }), 0), _vm._v(" "), _c("div", {
    staticClass: "bg-slate-50 rounded-2xl p-6 mb-8 border border-slate-100 font-bold"
  }, [_c("div", {
    staticClass: "flex justify-between mb-3 text-sm text-slate-500"
  }, [_c("span", [_vm._v("Tổng đơn")]), _c("span", [_vm._v(_vm._s(_vm.formatPrice((_vm$currentCapture6 = _vm.currentCapture) === null || _vm$currentCapture6 === void 0 ? void 0 : _vm$currentCapture6.total_amount)))])]), _vm._v(" "), _c("div", {
    staticClass: "pt-4 border-t border-slate-200 flex justify-between items-center"
  }, [_c("span", {
    staticClass: "text-lg"
  }, [_vm._v("Viettel Amount")]), _c("span", {
    staticClass: "text-2xl text-blue-600"
  }, [_vm._v(_vm._s(_vm.formatPrice((_vm$currentCapture7 = _vm.currentCapture) === null || _vm$currentCapture7 === void 0 ? void 0 : _vm$currentCapture7.viettel_amount)))])])]), _vm._v(" "), _vm._m(2)])]), _vm._v(" "), _c("el-dialog", {
    attrs: {
      title: "Order Receipt Generated",
      visible: _vm.captureVisible,
      width: "500px",
      "custom-class": "mobile-dialog"
    },
    on: {
      "update:visible": function updateVisible($event) {
        _vm.captureVisible = $event;
      }
    }
  }, [_c("div", {
    staticClass: "text-center"
  }, [_c("img", {
    staticClass: "w-full max-h-[70vh] w-auto mx-auto rounded-2xl border mb-6 shadow-sm",
    attrs: {
      src: _vm.captureImg
    }
  }), _c("el-button", {
    attrs: {
      type: "primary",
      size: "medium",
      icon: "el-icon-download"
    },
    on: {
      click: _vm.downloadCapture
    }
  }, [_vm._v("Download Image")])], 1)])], 1);
};
var staticRenderFns = [function () {
  var _vm = this,
    _c = _vm._self._c;
  return _c("div", {
    staticClass: "mobile-hide"
  }, [_c("h1", {
    staticClass: "text-3xl font-bold text-gray-900 tracking-tight"
  }, [_vm._v("Order Center")]), _vm._v(" "), _c("p", {
    staticClass: "text-slate-500 mt-1"
  }, [_vm._v("Manage sales, shipping, and receipt generation")])]);
}, function () {
  var _vm = this,
    _c = _vm._self._c;
  return _c("div", {
    staticClass: "text-center mb-6"
  }, [_c("h1", {
    staticClass: "text-4xl font-extrabold tracking-[4px] m-0 text-slate-900"
  }, [_vm._v("VINTAGE")]), _vm._v(" "), _c("p", {
    staticClass: "text-sm font-bold tracking-[2px] text-blue-600 mt-1 uppercase"
  }, [_vm._v("Always True")])]);
}, function () {
  var _vm = this,
    _c = _vm._self._c;
  return _c("div", {
    staticClass: "pt-6 border-t border-dashed border-slate-200 text-slate-500"
  }, [_c("div", {
    staticClass: "flex gap-4 mb-4"
  }, [_c("div", {
    staticClass: "flex-1"
  }, [_c("h4", {
    staticClass: "font-black text-slate-800 uppercase mb-2 text-[10px]"
  }, [_vm._v("🚚 GIAO HÀNG")]), _vm._v(" "), _c("ul", {
    staticClass: "list-none p-0 text-[10px] space-y-1"
  }, [_c("li", [_vm._v("• Viettel Post: HCM 1-2 ngày, tỉnh 3-5 ngày.")]), _vm._v(" "), _c("li", [_vm._v("• Phí ship đồng giá: 20k toàn quốc.")]), _vm._v(" "), _c("li", [_vm._v("• Hàng được giặt sấy, là ủi và xếp gọn tỉ mỉ.")])])]), _vm._v(" "), _c("div", {
    staticClass: "flex-1"
  }, [_c("h4", {
    staticClass: "font-black text-slate-800 uppercase mb-2 text-[10px]"
  }, [_vm._v("🔄 ĐỔI TRẢ")]), _vm._v(" "), _c("ul", {
    staticClass: "list-none p-0 text-[10px] space-y-1"
  }, [_c("li", [_vm._v("• Đổi trả nếu hàng không đúng mô tả, sai số đo >2cm.")]), _vm._v(" "), _c("li", [_vm._v("• Không nhận trả hàng với lý do cá nhân (không thích, không hợp).")])])])]), _vm._v(" "), _c("div", {
    staticClass: "text-center border-t border-dashed border-slate-200 pt-4"
  }, [_c("h3", {
    staticClass: "font-bold text-sm mb-1 text-slate-900"
  }, [_vm._v("Thank you for your order!")]), _vm._v(" "), _c("p", {
    staticClass: "text-[10px] text-slate-400"
  }, [_vm._v("Vintage Always True • Vintage Menswear")])])]);
}];
render._withStripped = true;


/***/ }),

/***/ "./node_modules/babel-loader/lib/index.js?!./node_modules/vue-loader/lib/loaders/templateLoader.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/ProductCreate.vue?vue&type=template&id=1e3b1bea":
/*!****************************************************************************************************************************************************************************************************************************************************!*\
  !*** ./node_modules/babel-loader/lib??ref--4-0!./node_modules/vue-loader/lib/loaders/templateLoader.js??ref--6!./node_modules/vue-loader/lib??vue-loader-options!./resources/js/components/salev2/ProductCreate.vue?vue&type=template&id=1e3b1bea ***!
  \****************************************************************************************************************************************************************************************************************************************************/
/*! exports provided: render, staticRenderFns */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "render", function() { return render; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "staticRenderFns", function() { return staticRenderFns; });
var render = function render() {
  var _vm = this,
    _c = _vm._self._c;
  return _c("div", {
    staticClass: "px-4 md:px-6 py-6 max-w-4xl mx-auto"
  }, [_c("div", {
    staticClass: "flex justify-between items-center mb-6"
  }, [_c("div", {
    staticClass: "flex items-center gap-3"
  }, [_c("el-button", {
    attrs: {
      icon: "el-icon-back",
      circle: ""
    },
    on: {
      click: function click($event) {
        return _vm.$router.push("/salev2/products");
      }
    }
  }), _vm._v(" "), _c("h1", {
    staticClass: "text-2xl font-bold text-slate-800"
  }, [_vm._v("Create New Product")])], 1)]), _vm._v(" "), _c("div", {
    staticClass: "bg-white rounded-xl shadow-sm border border-slate-100 overflow-hidden"
  }, [_c("div", {
    staticClass: "p-6 grid grid-cols-1 md:grid-cols-2 gap-8"
  }, [_c("div", {
    staticClass: "space-y-5"
  }, [_c("div", [_c("label", {
    staticClass: "block text-sm font-bold text-slate-700 mb-2"
  }, [_vm._v("Product Type")]), _vm._v(" "), _c("el-select", {
    staticClass: "w-full",
    attrs: {
      placeholder: "Select Type"
    },
    on: {
      change: _vm.handleTypeChange
    },
    model: {
      value: _vm.form.type,
      callback: function callback($$v) {
        _vm.$set(_vm.form, "type", $$v);
      },
      expression: "form.type"
    }
  }, _vm._l(_vm.types, function (type) {
    return _c("el-option", {
      key: type,
      attrs: {
        label: type,
        value: type
      }
    });
  }), 1)], 1), _vm._v(" "), _c("div", [_c("label", {
    staticClass: "block text-sm font-bold text-slate-700 mb-2"
  }, [_vm._v("\n                      Quick Input "), _c("span", {
    staticClass: "text-blue-500 font-normal ml-2 text-xs"
  }, [_vm._v("(" + _vm._s(_vm.sizeConfig.labels) + ")")])]), _vm._v(" "), _c("el-input", {
    attrs: {
      placeholder: "Paste measurements (e.g., 47 54 76 60)"
    },
    on: {
      input: _vm.parseQuickInput
    },
    model: {
      value: _vm.quickInput,
      callback: function callback($$v) {
        _vm.quickInput = $$v;
      },
      expression: "quickInput"
    }
  }, [_c("template", {
    slot: "prepend"
  }, [_c("i", {
    staticClass: "el-icon-magic-stick"
  })])], 2), _vm._v(" "), _c("p", {
    staticClass: "text-[10px] text-slate-400 mt-1"
  }, [_vm._v("Enter numbers separated by space to auto-fill sizes below.")])], 1), _vm._v(" "), _c("div", {
    staticClass: "bg-slate-50 rounded-lg p-5 border border-slate-200"
  }, [_c("h3", {
    staticClass: "font-bold text-slate-800 mb-4 border-b border-slate-200 pb-2"
  }, [_vm._v("Measurements")]), _vm._v(" "), _vm.form.type ? _c("div", {
    staticClass: "grid grid-cols-2 gap-4"
  }, _vm._l(_vm.activeSizes, function (size, idx) {
    return _c("div", {
      key: idx
    }, [_c("label", {
      staticClass: "block text-xs font-bold text-slate-500 mb-1 uppercase"
    }, [_vm._v(_vm._s(size))]), _vm._v(" "), _c("el-input", {
      attrs: {
        size: "small"
      },
      model: {
        value: _vm.form.sizes[idx].value,
        callback: function callback($$v) {
          _vm.$set(_vm.form.sizes[idx], "value", $$v);
        },
        expression: "form.sizes[idx].value"
      }
    }, [_c("template", {
      slot: "append"
    }, [_vm._v("cm")])], 2)], 1);
  }), 0) : _c("div", {
    staticClass: "text-center py-4 text-slate-400 text-xs"
  }, [_vm._v("Select Product Type first")])]), _vm._v(" "), _c("div", [_c("label", {
    staticClass: "block text-sm font-bold text-slate-700 mb-2"
  }, [_vm._v("Price")]), _vm._v(" "), _c("el-input-number", {
    staticClass: "w-full",
    attrs: {
      step: 50000
    },
    model: {
      value: _vm.form.price,
      callback: function callback($$v) {
        _vm.$set(_vm.form, "price", $$v);
      },
      expression: "form.price"
    }
  })], 1), _vm._v(" "), _c("div", [_c("label", {
    staticClass: "block text-sm font-bold text-slate-700 mb-2"
  }, [_vm._v("Description")]), _vm._v(" "), _c("el-input", {
    attrs: {
      type: "textarea",
      rows: 3,
      placeholder: "Product details..."
    },
    model: {
      value: _vm.form.description,
      callback: function callback($$v) {
        _vm.$set(_vm.form, "description", $$v);
      },
      expression: "form.description"
    }
  })], 1)]), _vm._v(" "), _c("div", [_c("label", {
    staticClass: "block text-sm font-bold text-slate-700 mb-2"
  }, [_vm._v("Thumbnail Image (Required)")]), _vm._v(" "), _c("div", {
    staticClass: "border-2 border-dashed border-slate-200 rounded-lg p-6 text-center hover:bg-slate-50 transition-colors relative cursor-pointer h-full min-h-[300px] flex flex-col items-center justify-center bg-white",
    on: {
      click: function click($event) {
        return _vm.$refs.fileInput.click();
      }
    }
  }, [_c("input", {
    ref: "fileInput",
    staticClass: "hidden",
    attrs: {
      type: "file",
      accept: "image/*"
    },
    on: {
      change: _vm.handleFileChange
    }
  }), _vm._v(" "), _vm.previewUrl ? _c("div", {
    staticClass: "w-full h-full flex items-center justify-center"
  }, [_c("img", {
    staticClass: "max-h-[400px] w-full object-contain",
    attrs: {
      src: _vm.previewUrl
    }
  })]) : _c("div", [_c("i", {
    staticClass: "el-icon-upload text-4xl text-slate-300 mb-3"
  }), _vm._v(" "), _c("p", {
    staticClass: "text-sm text-slate-500 font-bold"
  }, [_vm._v("Click to upload image")])])])])]), _vm._v(" "), _c("div", {
    staticClass: "px-6 py-4 bg-slate-50 border-t border-slate-100 flex justify-end gap-3"
  }, [_c("el-button", {
    on: {
      click: function click($event) {
        return _vm.$router.back();
      }
    }
  }, [_vm._v("Cancel")]), _vm._v(" "), _c("el-button", {
    attrs: {
      type: "primary",
      loading: _vm.loading
    },
    on: {
      click: _vm.submitByAjax
    }
  }, [_vm._v("Create Product")])], 1)])]);
};
var staticRenderFns = [];
render._withStripped = true;


/***/ }),

/***/ "./node_modules/babel-loader/lib/index.js?!./node_modules/vue-loader/lib/loaders/templateLoader.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/ProductEdit.vue?vue&type=template&id=635d4f78":
/*!**************************************************************************************************************************************************************************************************************************************************!*\
  !*** ./node_modules/babel-loader/lib??ref--4-0!./node_modules/vue-loader/lib/loaders/templateLoader.js??ref--6!./node_modules/vue-loader/lib??vue-loader-options!./resources/js/components/salev2/ProductEdit.vue?vue&type=template&id=635d4f78 ***!
  \**************************************************************************************************************************************************************************************************************************************************/
/*! exports provided: render, staticRenderFns */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "render", function() { return render; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "staticRenderFns", function() { return staticRenderFns; });
var render = function render() {
  var _vm = this,
    _c = _vm._self._c;
  return _c("div", {
    staticClass: "px-4 md:px-6 py-8"
  }, [_c("div", {
    staticClass: "flex justify-between items-center mb-6"
  }, [_vm._m(0), _vm._v(" "), _c("el-button", {
    attrs: {
      icon: "el-icon-back",
      circle: ""
    },
    on: {
      click: function click($event) {
        return _vm.$router.push("/salev2/products");
      }
    }
  })], 1), _vm._v(" "), _c("div", {
    directives: [{
      name: "loading",
      rawName: "v-loading",
      value: _vm.loading,
      expression: "loading"
    }],
    staticClass: "grid grid-cols-1 md:grid-cols-3 gap-6"
  }, [_c("div", {
    staticClass: "space-y-6"
  }, [_c("div", {
    staticClass: "bg-white p-6 rounded-2xl shadow-sm border border-slate-100"
  }, [_c("label", {
    staticClass: "block text-slate-400 font-bold text-xs uppercase mb-2 tracking-wider"
  }, [_vm._v("Product Type")]), _vm._v(" "), _c("el-select", {
    staticClass: "w-full",
    attrs: {
      placeholder: "Select Type",
      disabled: ""
    },
    model: {
      value: _vm.form.type,
      callback: function callback($$v) {
        _vm.$set(_vm.form, "type", $$v);
      },
      expression: "form.type"
    }
  }, _vm._l(_vm.productTypes, function (t) {
    return _c("el-option", {
      key: t,
      attrs: {
        label: t,
        value: t
      }
    });
  }), 1)], 1), _vm._v(" "), _c("div", {
    staticClass: "bg-white p-6 rounded-2xl shadow-sm border border-slate-100"
  }, [_c("label", {
    staticClass: "block text-slate-400 font-bold text-xs uppercase mb-2 tracking-wider"
  }, [_vm._v("Thumbnail Image")]), _vm._v(" "), _c("div", {
    staticClass: "relative group cursor-pointer",
    on: {
      click: function click($event) {
        return _vm.$refs.fileInput.click();
      }
    }
  }, [_vm.imagePreview ? _c("div", {
    staticClass: "aspect-[3/4] rounded-lg overflow-hidden bg-slate-50 border border-slate-200"
  }, [_c("img", {
    staticClass: "w-full h-full object-contain",
    attrs: {
      src: _vm.imagePreview
    }
  })]) : _c("div", {
    staticClass: "aspect-[3/4] rounded-lg bg-slate-50 border-2 border-dashed border-slate-200 flex flex-col items-center justify-center text-slate-400"
  }, [_c("i", {
    staticClass: "el-icon-picture text-3xl mb-2"
  }), _vm._v(" "), _c("span", {
    staticClass: "text-xs font-bold uppercase"
  }, [_vm._v("Upload Image")])]), _vm._v(" "), _vm._m(1)]), _vm._v(" "), _c("input", {
    ref: "fileInput",
    staticClass: "hidden",
    attrs: {
      type: "file",
      accept: "image/*"
    },
    on: {
      change: _vm.handleFileChange
    }
  })])]), _vm._v(" "), _c("div", {
    staticClass: "md:col-span-2 space-y-6"
  }, [_c("div", {
    staticClass: "bg-white p-6 rounded-2xl shadow-sm border border-slate-100"
  }, [_c("h2", {
    staticClass: "font-bold text-slate-800 text-lg mb-4"
  }, [_vm._v("Product Details")]), _vm._v(" "), _c("div", {
    staticClass: "space-y-4"
  }, [_c("div", [_c("label", {
    staticClass: "flex justify-between items-center text-slate-400 font-bold text-xs uppercase mb-2 tracking-wider"
  }, [_c("span", [_vm._v("Quick Input (Space separated)")]), _vm._v(" "), _c("span", {
    staticClass: "text-blue-600 cursor-pointer",
    on: {
      click: function click($event) {
        _vm.$message.info("Format: " + (_vm.activeMapping.labels || "Select type first"));
      }
    }
  }, [_c("i", {
    staticClass: "el-icon-info"
  }), _vm._v(" Help\n                           ")])]), _vm._v(" "), _c("el-input", {
    attrs: {
      placeholder: "e.g. 47 54 76 60"
    },
    on: {
      input: _vm.parseQuickInput
    },
    model: {
      value: _vm.quickInput,
      callback: function callback($$v) {
        _vm.quickInput = $$v;
      },
      expression: "quickInput"
    }
  }, [_c("template", {
    slot: "prepend"
  }, [_c("i", {
    staticClass: "el-icon-magic-stick"
  })])], 2), _vm._v(" "), _vm.activeMapping ? _c("p", {
    staticClass: "text-[10px] text-slate-400 mt-2"
  }, [_vm._v("\n                           Order: "), _c("span", {
    staticClass: "font-bold text-slate-600"
  }, [_vm._v(_vm._s(_vm.activeMapping.labels))])]) : _vm._e()], 1), _vm._v(" "), _c("div", {
    staticClass: "grid grid-cols-2 gap-4"
  }, [_c("div", [_c("label", {
    staticClass: "block text-slate-400 font-bold text-xs uppercase mb-2 tracking-wider"
  }, [_vm._v("Price")]), _vm._v(" "), _c("el-input-number", {
    staticClass: "w-full",
    attrs: {
      step: 10000,
      "controls-position": "right"
    },
    model: {
      value: _vm.form.price,
      callback: function callback($$v) {
        _vm.$set(_vm.form, "price", $$v);
      },
      expression: "form.price"
    }
  })], 1), _vm._v(" "), _c("div", [_c("label", {
    staticClass: "block text-slate-400 font-bold text-xs uppercase mb-2 tracking-wider"
  }, [_vm._v("Description")]), _vm._v(" "), _c("el-input", {
    attrs: {
      placeholder: "Optional notes..."
    },
    model: {
      value: _vm.form.description,
      callback: function callback($$v) {
        _vm.$set(_vm.form, "description", $$v);
      },
      expression: "form.description"
    }
  })], 1)])])]), _vm._v(" "), _vm.form.type ? _c("div", {
    staticClass: "bg-white p-6 rounded-2xl shadow-sm border border-slate-100"
  }, [_c("h2", {
    staticClass: "font-bold text-slate-800 text-lg mb-4"
  }, [_vm._v("Measurements")]), _vm._v(" "), _c("div", {
    staticClass: "grid grid-cols-2 md:grid-cols-3 gap-4"
  }, _vm._l(_vm.allSizes, function (sizeName, idx) {
    return _c("div", {
      directives: [{
        name: "show",
        rawName: "v-show",
        value: _vm.isSizeRelevant(idx),
        expression: "isSizeRelevant(idx)"
      }],
      key: sizeName,
      staticClass: "bg-slate-50 p-3 rounded-lg border border-slate-100"
    }, [_c("label", {
      staticClass: "block text-[10px] font-black uppercase text-slate-400 marginBottom-1"
    }, [_vm._v(_vm._s(sizeName))]), _vm._v(" "), _c("el-input", {
      attrs: {
        size: "medium",
        placeholder: "-"
      },
      model: {
        value: _vm.form.sizes[idx].value,
        callback: function callback($$v) {
          _vm.$set(_vm.form.sizes[idx], "value", $$v);
        },
        expression: "form.sizes[idx].value"
      }
    }, [_c("template", {
      slot: "append"
    }, [_vm._v("cm")])], 2)], 1);
  }), 0)]) : _vm._e(), _vm._v(" "), _c("div", {
    staticClass: "fixed bottom-0 md:relative left-0 right-0 p-4 md:p-0 bg-white md:bg-transparent border-t md:border-t-0 border-slate-200 z-40 flex justify-end gap-3"
  }, [_c("el-button", {
    on: {
      click: function click($event) {
        return _vm.$router.push("/salev2/products");
      }
    }
  }, [_vm._v("Cancel")]), _vm._v(" "), _c("el-button", {
    staticClass: "w-full md:w-auto px-8 shadow-lg shadow-blue-200",
    attrs: {
      type: "primary",
      loading: _vm.submitting
    },
    on: {
      click: _vm.submit
    }
  }, [_vm._v("\n                   Save Changes\n               ")])], 1)])])]);
};
var staticRenderFns = [function () {
  var _vm = this,
    _c = _vm._self._c;
  return _c("div", [_c("h1", {
    staticClass: "text-3xl font-bold text-gray-900 tracking-tight"
  }, [_vm._v("Edit Product")]), _vm._v(" "), _c("p", {
    staticClass: "text-slate-500 mt-1"
  }, [_vm._v("Update product details and measurements")])]);
}, function () {
  var _vm = this,
    _c = _vm._self._c;
  return _c("div", {
    staticClass: "absolute inset-0 bg-black/50 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center rounded-lg"
  }, [_c("span", {
    staticClass: "text-white font-bold uppercase tracking-widest text-xs border border-white px-4 py-2 rounded"
  }, [_vm._v("Change")])]);
}];
render._withStripped = true;


/***/ }),

/***/ "./node_modules/babel-loader/lib/index.js?!./node_modules/vue-loader/lib/loaders/templateLoader.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/ProductVault.vue?vue&type=template&id=1a1232b4":
/*!***************************************************************************************************************************************************************************************************************************************************!*\
  !*** ./node_modules/babel-loader/lib??ref--4-0!./node_modules/vue-loader/lib/loaders/templateLoader.js??ref--6!./node_modules/vue-loader/lib??vue-loader-options!./resources/js/components/salev2/ProductVault.vue?vue&type=template&id=1a1232b4 ***!
  \***************************************************************************************************************************************************************************************************************************************************/
/*! exports provided: render, staticRenderFns */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "render", function() { return render; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "staticRenderFns", function() { return staticRenderFns; });
var render = function render() {
  var _vm = this,
    _c = _vm._self._c;
  return _c("div", {
    staticClass: "px-4 md:px-6 py-8"
  }, [_c("div", {
    staticClass: "flex justify-between items-center mb-6 mobile-hide"
  }, [_vm._m(0), _vm._v(" "), _c("div", {
    staticClass: "flex gap-3"
  }, [_c("el-button", {
    attrs: {
      icon: "el-icon-plus",
      plain: ""
    },
    on: {
      click: function click($event) {
        return _vm.$router.push("/salev2/product/create");
      }
    }
  }, [_vm._v("Create")]), _vm._v(" "), _c("el-button", {
    attrs: {
      type: "primary",
      icon: "el-icon-upload"
    },
    on: {
      click: function click($event) {
        return _vm.$router.push("/salev2/upload");
      }
    }
  }, [_vm._v("Bulk Upload")])], 1)]), _vm._v(" "), _c("div", {
    staticClass: "mb-6"
  }, [_c("el-tabs", {
    staticClass: "premium-tabs",
    attrs: {
      type: "card"
    },
    on: {
      "tab-click": _vm.handleTabClick
    },
    model: {
      value: _vm.activeTab,
      callback: function callback($$v) {
        _vm.activeTab = $$v;
      },
      expression: "activeTab"
    }
  }, [_c("el-tab-pane", {
    attrs: {
      label: "Available",
      name: "AVAILABLE"
    }
  }), _vm._v(" "), _c("el-tab-pane", {
    attrs: {
      label: "All Products",
      name: "ALL"
    }
  })], 1)], 1), _vm._v(" "), _c("div", {
    staticClass: "bg-white p-4 md:p-6 rounded-2xl shadow-sm border border-slate-100 mb-8 space-y-4"
  }, [_c("div", {
    staticClass: "flex justify-between items-center md:hidden pb-2 mb-2 border-b border-slate-100 cursor-pointer"
  }, [_c("div", {
    staticClass: "flex items-center gap-2",
    on: {
      click: function click($event) {
        _vm.mobileFiltersOpen = !_vm.mobileFiltersOpen;
      }
    }
  }, [_c("i", {
    staticClass: "el-icon-s-operation text-blue-600"
  }), _vm._v(" "), _c("span", {
    staticClass: "font-bold text-slate-700 text-sm"
  }, [_vm._v("Filters")]), _vm._v(" "), _c("i", {
    staticClass: "text-slate-400 font-bold",
    "class": _vm.mobileFiltersOpen ? "el-icon-arrow-up" : "el-icon-arrow-down"
  })]), _vm._v(" "), _c("el-button", {
    attrs: {
      type: "primary",
      size: "mini",
      icon: "el-icon-plus",
      circle: ""
    },
    on: {
      click: function click($event) {
        return _vm.$router.push("/salev2/product/create");
      }
    }
  })], 1), _vm._v(" "), _c("div", {
    "class": _vm.mobileFiltersOpen ? "block" : "hidden md:block"
  }, [_c("div", {
    staticClass: "flex flex-col lg:flex-row gap-6 items-start"
  }, [_c("div", {
    staticClass: "w-full lg:w-64 flex-shrink-0 space-y-3"
  }, [_c("div", [_c("label", {
    staticClass: "text-[10px] font-black uppercase text-slate-400 mb-1 block tracking-wider"
  }, [_vm._v("Search")]), _vm._v(" "), _c("el-input", {
    staticClass: "w-full",
    attrs: {
      placeholder: "Name, code...",
      "prefix-icon": "el-icon-search",
      clearable: ""
    },
    on: {
      input: _vm.debouncedFetch
    },
    model: {
      value: _vm.search,
      callback: function callback($$v) {
        _vm.search = $$v;
      },
      expression: "search"
    }
  })], 1), _vm._v(" "), _c("div", [_c("label", {
    staticClass: "text-[10px] font-black uppercase text-slate-400 mb-1 block tracking-wider"
  }, [_vm._v("Product Type")]), _vm._v(" "), _c("el-select", {
    staticClass: "w-full",
    attrs: {
      placeholder: "Any Type",
      clearable: ""
    },
    on: {
      change: function change($event) {
        return _vm.fetchProducts(1);
      }
    },
    model: {
      value: _vm.filterType,
      callback: function callback($$v) {
        _vm.filterType = $$v;
      },
      expression: "filterType"
    }
  }, _vm._l(_vm.productTypes, function (t) {
    return _c("el-option", {
      key: t,
      attrs: {
        label: t,
        value: t
      }
    });
  }), 1)], 1)]), _vm._v(" "), _c("div", {
    staticClass: "flex-1 w-full"
  }, [_c("div", {
    staticClass: "flex flex-col md:flex-row flex-wrap items-end gap-3"
  }, [_vm._l(_vm.sizeFilters, function (filter, index) {
    return _c("div", {
      key: index,
      staticClass: "flex flex-wrap md:flex-nowrap items-end gap-2 p-3 bg-slate-50/50 rounded-xl border border-dashed border-slate-200"
    }, [_c("div", {
      staticClass: "w-[calc(100%-2.5rem)] md:w-32 order-1 md:order-1"
    }, [_c("label", {
      staticClass: "text-[10px] font-black uppercase text-slate-400 mb-1 block tracking-wider"
    }, [_vm._v("Key")]), _vm._v(" "), _c("el-select", {
      staticClass: "w-full",
      attrs: {
        placeholder: "Key",
        clearable: "",
        size: "mini"
      },
      model: {
        value: filter.key,
        callback: function callback($$v) {
          _vm.$set(filter, "key", $$v);
        },
        expression: "filter.key"
      }
    }, _vm._l(_vm.sizes, function (s) {
      return _c("el-option", {
        key: s,
        attrs: {
          label: s,
          value: s
        }
      });
    }), 1)], 1), _vm._v(" "), _c("div", {
      staticClass: "w-[48%] md:w-28 order-3 md:order-2"
    }, [_c("label", {
      staticClass: "text-[10px] font-black uppercase text-slate-400 mb-1 block tracking-wider"
    }, [_vm._v("From")]), _vm._v(" "), _c("el-input-number", {
      staticClass: "w-full",
      attrs: {
        precision: 1,
        step: 0.5,
        placeholder: "From",
        "controls-position": "right",
        size: "mini"
      },
      model: {
        value: filter.from,
        callback: function callback($$v) {
          _vm.$set(filter, "from", $$v);
        },
        expression: "filter.from"
      }
    })], 1), _vm._v(" "), _c("div", {
      staticClass: "w-[48%] md:w-28 order-4 md:order-3"
    }, [_c("label", {
      staticClass: "text-[10px] font-black uppercase text-slate-400 mb-1 block tracking-wider"
    }, [_vm._v("To")]), _vm._v(" "), _c("el-input-number", {
      staticClass: "w-full",
      attrs: {
        precision: 1,
        step: 0.5,
        placeholder: "To",
        "controls-position": "right",
        size: "mini"
      },
      model: {
        value: filter.to,
        callback: function callback($$v) {
          _vm.$set(filter, "to", $$v);
        },
        expression: "filter.to"
      }
    })], 1), _vm._v(" "), _c("div", {
      staticClass: "w-8 md:w-auto flex justify-end md:block mb-0.5 md:mb-0 order-2 md:order-4"
    }, [_vm.sizeFilters.length > 1 ? _c("el-button", {
      attrs: {
        type: "danger",
        icon: "el-icon-minus",
        circle: "",
        size: "mini",
        plain: ""
      },
      on: {
        click: function click($event) {
          return _vm.removeSizeFilter(index);
        }
      }
    }) : _c("div", {
      staticClass: "w-7 h-7 hidden md:block"
    })], 1)]);
  }), _vm._v(" "), _c("div", {
    staticClass: "pb-3 md:pb-3.5"
  }, [_vm.sizeFilters.length < 3 ? _c("el-button", {
    attrs: {
      type: "primary",
      icon: "el-icon-plus",
      circle: "",
      size: "small",
      plain: ""
    },
    on: {
      click: _vm.addSizeFilter
    }
  }) : _vm._e()], 1)], 2)])]), _vm._v(" "), _c("div", {
    staticClass: "flex flex-col md:flex-row justify-between items-center gap-4 pt-4 border-t border-slate-100 mt-2"
  }, [_c("div", {
    staticClass: "flex flex-nowrap w-full md:w-auto gap-2"
  }, [_c("el-button", {
    staticClass: "flex-1 md:flex-none px-2",
    attrs: {
      size: "small",
      icon: "el-icon-check",
      plain: ""
    },
    on: {
      click: _vm.selectAllPage
    }
  }, [_vm._v("Select Page")]), _vm._v(" "), _c("el-button", {
    staticClass: "flex-1 md:flex-none px-2",
    attrs: {
      size: "small",
      type: "success",
      disabled: !_vm.selectedProducts.length,
      icon: "el-icon-share",
      plain: ""
    },
    on: {
      click: _vm.shareSelected
    }
  }, [_vm._v("\n                  Share "), _vm.selectedProducts.length ? _c("span", [_vm._v("(" + _vm._s(_vm.selectedProducts.length) + ")")]) : _vm._e()]), _vm._v(" "), _vm.selectedProducts.length ? _c("el-button", {
    staticClass: "flex-none px-3",
    attrs: {
      size: "small",
      type: "danger",
      plain: "",
      icon: "el-icon-delete"
    },
    on: {
      click: _vm.clearSelection
    }
  }, [_vm._v("Clear")]) : _vm._e()], 1), _vm._v(" "), _c("div", {
    staticClass: "flex gap-3 w-full md:w-auto"
  }, [_c("el-button", {
    staticClass: "flex-1 md:flex-none",
    attrs: {
      size: "medium"
    },
    on: {
      click: _vm.resetFilters
    }
  }, [_vm._v("Reset")]), _vm._v(" "), _c("el-button", {
    staticClass: "flex-1 md:flex-none px-6 shadow-md shadow-blue-200",
    attrs: {
      type: "primary",
      size: "medium"
    },
    on: {
      click: function click($event) {
        return _vm.fetchProducts(1);
      }
    }
  }, [_vm._v("Apply Filters")])], 1)])])]), _vm._v(" "), _c("div", {
    directives: [{
      name: "loading",
      rawName: "v-loading",
      value: _vm.loading,
      expression: "loading"
    }]
  }, [_c("div", {
    staticClass: "grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-6 gap-3"
  }, _vm._l(_vm.products, function (product) {
    return _c("div", {
      key: product.id,
      staticClass: "group relative bg-white rounded-xl overflow-hidden border border-slate-100 shadow-sm hover:shadow-xl transition-all duration-300",
      "class": {
        "ring-2 ring-blue-500": _vm.isSelected(product.id)
      }
    }, [_c("div", {
      staticClass: "aspect-[3/4] overflow-hidden bg-slate-50 relative group"
    }, [product.path_thumb ? _c("img", {
      staticClass: "w-full h-full object-contain hover:scale-105 transition-transform duration-500 p-1",
      attrs: {
        src: product.path_thumb
      }
    }) : _c("div", {
      staticClass: "w-full h-full flex items-center justify-center text-slate-300"
    }, [_c("i", {
      staticClass: "el-icon-picture text-4xl"
    })]), _vm._v(" "), _c("div", {
      staticClass: "absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex flex-col justify-center items-center gap-3 p-4"
    }, [_c("el-button", {
      staticClass: "w-28 shadow-lg transform translate-y-4 group-hover:translate-y-0 transition-transform duration-300",
      attrs: {
        type: "primary",
        size: "small"
      },
      on: {
        click: function click($event) {
          return _vm.editProduct(product);
        }
      }
    }, [_vm._v("Edit")]), _vm._v(" "), _c("el-button", {
      staticClass: "w-28 shadow-lg transform translate-y-4 group-hover:translate-y-0 transition-transform duration-300 delay-75",
      attrs: {
        type: "success",
        size: "small",
        icon: "el-icon-document-copy"
      },
      on: {
        click: function click($event) {
          return _vm.copyProductInfo(product);
        }
      }
    }, [_vm._v("Copy")])], 1)]), _vm._v(" "), _c("div", {
      staticClass: "p-3"
    }, [_c("div", {
      staticClass: "text-center mb-2"
    }, [_c("h4", {
      staticClass: "font-bold text-slate-800 text-sm uppercase leading-tight truncate"
    }, [_vm._v("\n                  " + _vm._s(product.name) + "\n              ")]), _vm._v(" "), _c("p", {
      staticClass: "text-blue-600 font-bold text-sm"
    }, [_vm._v(_vm._s(_vm.formatPrice(product.price)))])]), _vm._v(" "), _c("div", {
      staticClass: "text-center mb-3 px-1"
    }, [_c("span", {
      staticClass: "text-[10px] font-black bg-blue-50 text-blue-600 px-2 py-1 rounded leading-tight uppercase inline-block"
    }, [_vm._v("\n                " + _vm._s(_vm.formatSizes(product).join(" - ")) + "\n              ")])]), _vm._v(" "), _c("div", {
      staticClass: "flex justify-between items-center pt-2 border-t border-slate-50"
    }, [_c("el-tag", {
      staticClass: "text-[8px] font-bold tracking-tighter uppercase rounded-md px-1",
      attrs: {
        type: _vm.getStatusTagType(product.status),
        size: "mini"
      }
    }, [_vm._v("\n                  " + _vm._s(product.status) + "\n              ")]), _vm._v(" "), _c("el-checkbox", {
      staticClass: "transform scale-125",
      attrs: {
        value: _vm.isSelected(product.id)
      },
      on: {
        change: function change($event) {
          return _vm.toggleSelection(product);
        }
      }
    })], 1)])]);
  }), 0), _vm._v(" "), !_vm.products.length && !_vm.loading ? _c("div", {
    staticClass: "text-center py-20 bg-white rounded-3xl border border-dashed border-slate-200"
  }, [_c("i", {
    staticClass: "el-icon-search text-4xl text-slate-200 mb-4 block"
  }), _vm._v(" "), _c("p", {
    staticClass: "text-slate-400 font-medium"
  }, [_vm._v("No products found")])]) : _vm._e(), _vm._v(" "), _c("div", {
    staticClass: "mt-8 flex justify-center"
  }, [_c("el-pagination", {
    attrs: {
      background: "",
      layout: "prev, pager, next",
      "current-page": _vm.page,
      total: _vm.total,
      "page-size": 60
    },
    on: {
      "update:currentPage": function updateCurrentPage($event) {
        _vm.page = $event;
      },
      "update:current-page": function updateCurrentPage($event) {
        _vm.page = $event;
      },
      "current-change": _vm.handlePageChange
    }
  })], 1)])]);
};
var staticRenderFns = [function () {
  var _vm = this,
    _c = _vm._self._c;
  return _c("div", [_c("h1", {
    staticClass: "text-3xl font-bold text-gray-900 tracking-tight"
  }, [_vm._v("Product Vault")]), _vm._v(" "), _c("p", {
    staticClass: "text-slate-500 mt-1"
  }, [_vm._v("Browse and share available inventory")])]);
}];
render._withStripped = true;


/***/ }),

/***/ "./node_modules/babel-loader/lib/index.js?!./node_modules/vue-loader/lib/loaders/templateLoader.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/SmartUpload.vue?vue&type=template&id=30d6d489":
/*!**************************************************************************************************************************************************************************************************************************************************!*\
  !*** ./node_modules/babel-loader/lib??ref--4-0!./node_modules/vue-loader/lib/loaders/templateLoader.js??ref--6!./node_modules/vue-loader/lib??vue-loader-options!./resources/js/components/salev2/SmartUpload.vue?vue&type=template&id=30d6d489 ***!
  \**************************************************************************************************************************************************************************************************************************************************/
/*! exports provided: render, staticRenderFns */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "render", function() { return render; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "staticRenderFns", function() { return staticRenderFns; });
var render = function render() {
  var _vm = this,
    _c = _vm._self._c;
  return _c("div", {
    staticClass: "px-6 py-8"
  }, [_vm._m(0), _vm._v(" "), _c("div", {
    staticClass: "grid grid-cols-1 lg:grid-cols-2 gap-8 mb-8"
  }, [_c("div", {
    staticClass: "bg-white rounded-2xl shadow-sm border border-slate-100 p-6"
  }, [_c("div", {
    staticClass: "space-y-6"
  }, [_c("div", [_c("label", {
    staticClass: "text-xs font-black text-slate-400 uppercase tracking-widest block mb-2"
  }, [_vm._v("Product Type")]), _vm._v(" "), _c("el-select", {
    staticClass: "w-full shadow-sm",
    attrs: {
      placeholder: "Select Type"
    },
    model: {
      value: _vm.productType,
      callback: function callback($$v) {
        _vm.productType = $$v;
      },
      expression: "productType"
    }
  }, _vm._l(_vm.productTypes, function (type) {
    return _c("el-option", {
      key: type,
      attrs: {
        label: type,
        value: type
      }
    });
  }), 1)], 1), _vm._v(" "), _c("div", [_c("label", {
    staticClass: "text-xs font-black text-slate-400 uppercase tracking-widest block mb-2"
  }, [_vm._v("Select Images (Multiple)")]), _vm._v(" "), _c("div", {
    staticClass: "flex items-center gap-4"
  }, [_c("el-button", {
    attrs: {
      type: "info",
      plain: "",
      icon: "el-icon-files"
    },
    on: {
      click: _vm.triggerFileSelect
    }
  }, [_vm._v("\n              " + _vm._s(_vm.selectedFiles.length ? "".concat(_vm.selectedFiles.length, " files selected") : "Choose Files") + "\n            ")]), _vm._v(" "), _c("input", {
    ref: "fileInput",
    attrs: {
      type: "file",
      multiple: "",
      hidden: "",
      accept: "image/*"
    },
    on: {
      change: _vm.handleFileChange
    }
  }), _vm._v(" "), _c("el-button", {
    attrs: {
      type: "primary",
      disabled: !_vm.selectedFiles.length || _vm.uploading,
      loading: _vm.uploading,
      icon: "el-icon-upload"
    },
    on: {
      click: _vm.startUpload
    }
  }, [_vm._v("\n              Start Uploading\n            ")])], 1)])])]), _vm._v(" "), _c("div", {
    staticClass: "bg-slate-900 rounded-2xl p-6 shadow-xl text-white relative overflow-hidden"
  }, [_c("div", {
    staticClass: "relative z-10 flex flex-col h-full justify-between"
  }, [_c("div", [_c("p", {
    staticClass: "text-[10px] font-black opacity-40 uppercase tracking-[4px] mb-4"
  }, [_vm._v("Upload Progress")]), _vm._v(" "), _c("div", {
    staticClass: "flex items-end gap-2 mb-2"
  }, [_c("span", {
    staticClass: "text-4xl font-black"
  }, [_vm._v(_vm._s(_vm.progress) + "%")]), _vm._v(" "), _c("span", {
    staticClass: "text-xs opacity-60 mb-1"
  }, [_vm._v("Session: " + _vm._s(_vm.sessionId || "---"))])]), _vm._v(" "), _c("el-progress", {
    staticClass: "mb-6",
    attrs: {
      percentage: _vm.progress,
      "show-text": false,
      color: "#3b82f6"
    }
  })], 1), _vm._v(" "), _c("div", {
    staticClass: "grid grid-cols-4 gap-4 text-center"
  }, [_c("div", [_c("p", {
    staticClass: "text-xl font-black text-blue-400"
  }, [_vm._v(_vm._s(_vm.totalCount))]), _c("p", {
    staticClass: "text-[10px] opacity-40 font-bold uppercase"
  }, [_vm._v("Total")])]), _vm._v(" "), _c("div", [_c("p", {
    staticClass: "text-xl font-black text-amber-400"
  }, [_vm._v(_vm._s(_vm.processingCount))]), _c("p", {
    staticClass: "text-[10px] opacity-40 font-bold uppercase"
  }, [_vm._v("Proc")])]), _vm._v(" "), _c("div", [_c("p", {
    staticClass: "text-xl font-black text-emerald-400"
  }, [_vm._v(_vm._s(_vm.successCount))]), _c("p", {
    staticClass: "text-[10px] opacity-40 font-bold uppercase"
  }, [_vm._v("Done")])]), _vm._v(" "), _c("div", [_c("p", {
    staticClass: "text-xl font-black text-rose-400"
  }, [_vm._v(_vm._s(_vm.failedCount))]), _c("p", {
    staticClass: "text-[10px] opacity-40 font-bold uppercase"
  }, [_vm._v("Fail")])])])])])]), _vm._v(" "), _c("div", {
    staticClass: "grid grid-cols-1 xl:grid-cols-2 gap-8"
  }, [_c("div", {
    staticClass: "bg-white rounded-2xl shadow-sm border border-slate-100 overflow-hidden"
  }, [_c("div", {
    staticClass: "bg-slate-50 px-4 py-3 border-b border-slate-100 flex justify-between items-center"
  }, [_c("h3", {
    staticClass: "text-xs font-black text-slate-400 uppercase tracking-widest"
  }, [_vm._v("Network Upload Log")]), _vm._v(" "), _c("span", {
    staticClass: "text-[10px] bg-blue-100 text-blue-600 px-2 py-0.5 rounded-full font-bold"
  }, [_vm._v(_vm._s(_vm.uploadLog.length) + " Files")])]), _vm._v(" "), _c("div", {
    staticClass: "max-h-[500px] overflow-y-auto"
  }, [_c("el-table", {
    attrs: {
      data: _vm.uploadLog,
      size: "mini",
      "row-class-name": _vm.logRowClass
    }
  }, [_c("el-table-column", {
    attrs: {
      prop: "index",
      label: "#",
      width: "50",
      align: "center"
    }
  }), _vm._v(" "), _c("el-table-column", {
    attrs: {
      prop: "filename",
      label: "Filename",
      "min-width": "150",
      "show-overflow-tooltip": ""
    }
  }), _vm._v(" "), _c("el-table-column", {
    attrs: {
      label: "Status",
      width: "100",
      align: "center"
    },
    scopedSlots: _vm._u([{
      key: "default",
      fn: function fn(scope) {
        return [_c("el-tag", {
          staticClass: "uppercase font-black text-[8px] rounded-pill",
          attrs: {
            type: scope.row.status === "success" ? "success" : scope.row.status === "uploading" ? "info" : "danger",
            size: "mini",
            effect: "dark"
          }
        }, [_vm._v("\n                          " + _vm._s(scope.row.status) + "\n                      ")])];
      }
    }])
  }), _vm._v(" "), _c("el-table-column", {
    attrs: {
      prop: "message",
      label: "Message",
      "min-width": "200",
      "show-overflow-tooltip": ""
    }
  })], 1)], 1)]), _vm._v(" "), _c("div", {
    staticClass: "bg-white rounded-2xl shadow-sm border border-slate-100 overflow-hidden"
  }, [_c("div", {
    staticClass: "bg-slate-50 px-4 py-3 border-b border-slate-100 flex justify-between items-center"
  }, [_c("h3", {
    staticClass: "text-xs font-black text-slate-400 uppercase tracking-widest"
  }, [_vm._v("Backend Processing Log")]), _vm._v(" "), _vm.processingCount > 0 ? _c("span", {
    staticClass: "animate-pulse text-[10px] bg-amber-100 text-amber-600 px-2 py-0.5 rounded-full font-bold"
  }, [_vm._v("Processing...")]) : _vm._e()]), _vm._v(" "), _c("div", {
    staticClass: "max-h-[500px] overflow-y-auto"
  }, [_c("el-table", {
    attrs: {
      data: _vm.processingLog,
      size: "mini"
    }
  }, [_c("el-table-column", {
    attrs: {
      type: "index",
      label: "#",
      width: "50",
      align: "center"
    }
  }), _vm._v(" "), _c("el-table-column", {
    attrs: {
      prop: "file",
      label: "File",
      "min-width": "150",
      "show-overflow-tooltip": ""
    }
  }), _vm._v(" "), _c("el-table-column", {
    attrs: {
      label: "Status",
      width: "100",
      align: "center"
    },
    scopedSlots: _vm._u([{
      key: "default",
      fn: function fn(scope) {
        return [_c("el-tag", {
          staticClass: "uppercase font-black text-[8px] rounded-pill",
          attrs: {
            type: _vm.getProcStatusType(scope.row.status),
            size: "mini",
            effect: "dark"
          }
        }, [_vm._v("\n                          " + _vm._s(scope.row.status) + "\n                      ")])];
      }
    }])
  }), _vm._v(" "), _c("el-table-column", {
    attrs: {
      prop: "code",
      label: "Code",
      width: "100",
      align: "center"
    }
  }), _vm._v(" "), _c("el-table-column", {
    attrs: {
      prop: "reason",
      label: "Result/Error",
      "min-width": "200",
      "show-overflow-tooltip": ""
    }
  })], 1)], 1)])])]);
};
var staticRenderFns = [function () {
  var _vm = this,
    _c = _vm._self._c;
  return _c("div", {
    staticClass: "mb-8 mobile-hide"
  }, [_c("h1", {
    staticClass: "text-3xl font-bold text-gray-900 tracking-tight mb-2"
  }, [_vm._v("Batch Image Upload")]), _vm._v(" "), _c("p", {
    staticClass: "text-slate-500 italic"
  }, [_vm._v("Queue-based processing for massive product catalogs")])]);
}];
render._withStripped = true;


/***/ }),

/***/ "./node_modules/css-loader/index.js?!./node_modules/vue-loader/lib/loaders/stylePostLoader.js!./node_modules/postcss-loader/src/index.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/Dashboard.vue?vue&type=style&index=0&id=16392a9a&lang=css":
/*!*********************************************************************************************************************************************************************************************************************************************************************************************!*\
  !*** ./node_modules/css-loader??ref--5-1!./node_modules/vue-loader/lib/loaders/stylePostLoader.js!./node_modules/postcss-loader/src??ref--5-2!./node_modules/vue-loader/lib??vue-loader-options!./resources/js/components/salev2/Dashboard.vue?vue&type=style&index=0&id=16392a9a&lang=css ***!
  \*********************************************************************************************************************************************************************************************************************************************************************************************/
/*! no static exports found */
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__(/*! ../../../../node_modules/css-loader/lib/css-base.js */ "./node_modules/css-loader/lib/css-base.js")(false);
// imports


// module
exports.push([module.i, "\n.legacy-dash-table .el-table__header {\n    background-color: #f8fafc !important;\n}\n.legacy-dash-table .el-table__row {\n    font-size: 11px;\n}\n.legacy-dash-table .el-table__header th {\n    font-size: 10px;\n    font-weight: 900;\n    text-transform: uppercase;\n    color: #94a3b8;\n    padding: 8px 0;\n}\n", ""]);

// exports


/***/ }),

/***/ "./node_modules/css-loader/index.js?!./node_modules/vue-loader/lib/loaders/stylePostLoader.js!./node_modules/postcss-loader/src/index.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/OrderManager.vue?vue&type=style&index=0&id=64110eb0&lang=css":
/*!************************************************************************************************************************************************************************************************************************************************************************************************!*\
  !*** ./node_modules/css-loader??ref--5-1!./node_modules/vue-loader/lib/loaders/stylePostLoader.js!./node_modules/postcss-loader/src??ref--5-2!./node_modules/vue-loader/lib??vue-loader-options!./resources/js/components/salev2/OrderManager.vue?vue&type=style&index=0&id=64110eb0&lang=css ***!
  \************************************************************************************************************************************************************************************************************************************************************************************************/
/*! no static exports found */
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__(/*! ../../../../node_modules/css-loader/lib/css-base.js */ "./node_modules/css-loader/lib/css-base.js")(false);
// imports


// module
exports.push([module.i, "\n.scrollbar-hide::-webkit-scrollbar { display: none;\n}\n.scrollbar-hide { -ms-overflow-style: none; scrollbar-width: none;\n}\n@media (max-width: 768px) {\n.mobile-dialog { width: 95% !important; margin-top: 5vh !important;\n}\n.mobile-dialog .el-dialog__body { padding: 10px 15px;\n}\n}\n", ""]);

// exports


/***/ }),

/***/ "./node_modules/css-loader/index.js?!./node_modules/vue-loader/lib/loaders/stylePostLoader.js!./node_modules/postcss-loader/src/index.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/ProductVault.vue?vue&type=style&index=0&id=1a1232b4&lang=css":
/*!************************************************************************************************************************************************************************************************************************************************************************************************!*\
  !*** ./node_modules/css-loader??ref--5-1!./node_modules/vue-loader/lib/loaders/stylePostLoader.js!./node_modules/postcss-loader/src??ref--5-2!./node_modules/vue-loader/lib??vue-loader-options!./resources/js/components/salev2/ProductVault.vue?vue&type=style&index=0&id=1a1232b4&lang=css ***!
  \************************************************************************************************************************************************************************************************************************************************************************************************/
/*! no static exports found */
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__(/*! ../../../../node_modules/css-loader/lib/css-base.js */ "./node_modules/css-loader/lib/css-base.js")(false);
// imports


// module
exports.push([module.i, "\n.premium-tabs .el-tabs__item {\n    font-weight: 900;\n    text-transform: uppercase;\n    font-size: 12px;\n    letter-spacing: 1px;\n}\n.premium-tabs .el-tabs__header {\n    margin-bottom: 0;\n}\n", ""]);

// exports


/***/ }),

/***/ "./node_modules/css-loader/index.js?!./node_modules/vue-loader/lib/loaders/stylePostLoader.js!./node_modules/postcss-loader/src/index.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/SmartUpload.vue?vue&type=style&index=0&id=30d6d489&lang=css":
/*!***********************************************************************************************************************************************************************************************************************************************************************************************!*\
  !*** ./node_modules/css-loader??ref--5-1!./node_modules/vue-loader/lib/loaders/stylePostLoader.js!./node_modules/postcss-loader/src??ref--5-2!./node_modules/vue-loader/lib??vue-loader-options!./resources/js/components/salev2/SmartUpload.vue?vue&type=style&index=0&id=30d6d489&lang=css ***!
  \***********************************************************************************************************************************************************************************************************************************************************************************************/
/*! no static exports found */
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__(/*! ../../../../node_modules/css-loader/lib/css-base.js */ "./node_modules/css-loader/lib/css-base.js")(false);
// imports


// module
exports.push([module.i, "\n.el-table .success-row { background: #f0f9eb;\n}\n.el-table .danger-row { background: #fef0f0;\n}\n", ""]);

// exports


/***/ }),

/***/ "./node_modules/css-loader/lib/css-base.js":
/*!*************************************************!*\
  !*** ./node_modules/css-loader/lib/css-base.js ***!
  \*************************************************/
/*! no static exports found */
/***/ (function(module, exports) {

/*
	MIT License http://www.opensource.org/licenses/mit-license.php
	Author Tobias Koppers @sokra
*/
// css base code, injected by the css-loader
module.exports = function(useSourceMap) {
	var list = [];

	// return the list of modules as css string
	list.toString = function toString() {
		return this.map(function (item) {
			var content = cssWithMappingToString(item, useSourceMap);
			if(item[2]) {
				return "@media " + item[2] + "{" + content + "}";
			} else {
				return content;
			}
		}).join("");
	};

	// import a list of modules into the list
	list.i = function(modules, mediaQuery) {
		if(typeof modules === "string")
			modules = [[null, modules, ""]];
		var alreadyImportedModules = {};
		for(var i = 0; i < this.length; i++) {
			var id = this[i][0];
			if(typeof id === "number")
				alreadyImportedModules[id] = true;
		}
		for(i = 0; i < modules.length; i++) {
			var item = modules[i];
			// skip already imported module
			// this implementation is not 100% perfect for weird media query combinations
			//  when a module is imported multiple times with different media queries.
			//  I hope this will never occur (Hey this way we have smaller bundles)
			if(typeof item[0] !== "number" || !alreadyImportedModules[item[0]]) {
				if(mediaQuery && !item[2]) {
					item[2] = mediaQuery;
				} else if(mediaQuery) {
					item[2] = "(" + item[2] + ") and (" + mediaQuery + ")";
				}
				list.push(item);
			}
		}
	};
	return list;
};

function cssWithMappingToString(item, useSourceMap) {
	var content = item[1] || '';
	var cssMapping = item[3];
	if (!cssMapping) {
		return content;
	}

	if (useSourceMap && typeof btoa === 'function') {
		var sourceMapping = toComment(cssMapping);
		var sourceURLs = cssMapping.sources.map(function (source) {
			return '/*# sourceURL=' + cssMapping.sourceRoot + source + ' */'
		});

		return [content].concat(sourceURLs).concat([sourceMapping]).join('\n');
	}

	return [content].join('\n');
}

// Adapted from convert-source-map (MIT)
function toComment(sourceMap) {
	// eslint-disable-next-line no-undef
	var base64 = btoa(unescape(encodeURIComponent(JSON.stringify(sourceMap))));
	var data = 'sourceMappingURL=data:application/json;charset=utf-8;base64,' + base64;

	return '/*# ' + data + ' */';
}


/***/ }),

/***/ "./node_modules/style-loader/index.js!./node_modules/css-loader/index.js?!./node_modules/vue-loader/lib/loaders/stylePostLoader.js!./node_modules/postcss-loader/src/index.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/Dashboard.vue?vue&type=style&index=0&id=16392a9a&lang=css":
/*!*************************************************************************************************************************************************************************************************************************************************************************************************************************!*\
  !*** ./node_modules/style-loader!./node_modules/css-loader??ref--5-1!./node_modules/vue-loader/lib/loaders/stylePostLoader.js!./node_modules/postcss-loader/src??ref--5-2!./node_modules/vue-loader/lib??vue-loader-options!./resources/js/components/salev2/Dashboard.vue?vue&type=style&index=0&id=16392a9a&lang=css ***!
  \*************************************************************************************************************************************************************************************************************************************************************************************************************************/
/*! no static exports found */
/***/ (function(module, exports, __webpack_require__) {


var content = __webpack_require__(/*! !../../../../node_modules/css-loader??ref--5-1!../../../../node_modules/vue-loader/lib/loaders/stylePostLoader.js!../../../../node_modules/postcss-loader/src??ref--5-2!../../../../node_modules/vue-loader/lib??vue-loader-options!./Dashboard.vue?vue&type=style&index=0&id=16392a9a&lang=css */ "./node_modules/css-loader/index.js?!./node_modules/vue-loader/lib/loaders/stylePostLoader.js!./node_modules/postcss-loader/src/index.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/Dashboard.vue?vue&type=style&index=0&id=16392a9a&lang=css");

if(typeof content === 'string') content = [[module.i, content, '']];

var transform;
var insertInto;



var options = {"hmr":true}

options.transform = transform
options.insertInto = undefined;

var update = __webpack_require__(/*! ../../../../node_modules/style-loader/lib/addStyles.js */ "./node_modules/style-loader/lib/addStyles.js")(content, options);

if(content.locals) module.exports = content.locals;

if(false) {}

/***/ }),

/***/ "./node_modules/style-loader/index.js!./node_modules/css-loader/index.js?!./node_modules/vue-loader/lib/loaders/stylePostLoader.js!./node_modules/postcss-loader/src/index.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/OrderManager.vue?vue&type=style&index=0&id=64110eb0&lang=css":
/*!****************************************************************************************************************************************************************************************************************************************************************************************************************************!*\
  !*** ./node_modules/style-loader!./node_modules/css-loader??ref--5-1!./node_modules/vue-loader/lib/loaders/stylePostLoader.js!./node_modules/postcss-loader/src??ref--5-2!./node_modules/vue-loader/lib??vue-loader-options!./resources/js/components/salev2/OrderManager.vue?vue&type=style&index=0&id=64110eb0&lang=css ***!
  \****************************************************************************************************************************************************************************************************************************************************************************************************************************/
/*! no static exports found */
/***/ (function(module, exports, __webpack_require__) {


var content = __webpack_require__(/*! !../../../../node_modules/css-loader??ref--5-1!../../../../node_modules/vue-loader/lib/loaders/stylePostLoader.js!../../../../node_modules/postcss-loader/src??ref--5-2!../../../../node_modules/vue-loader/lib??vue-loader-options!./OrderManager.vue?vue&type=style&index=0&id=64110eb0&lang=css */ "./node_modules/css-loader/index.js?!./node_modules/vue-loader/lib/loaders/stylePostLoader.js!./node_modules/postcss-loader/src/index.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/OrderManager.vue?vue&type=style&index=0&id=64110eb0&lang=css");

if(typeof content === 'string') content = [[module.i, content, '']];

var transform;
var insertInto;



var options = {"hmr":true}

options.transform = transform
options.insertInto = undefined;

var update = __webpack_require__(/*! ../../../../node_modules/style-loader/lib/addStyles.js */ "./node_modules/style-loader/lib/addStyles.js")(content, options);

if(content.locals) module.exports = content.locals;

if(false) {}

/***/ }),

/***/ "./node_modules/style-loader/index.js!./node_modules/css-loader/index.js?!./node_modules/vue-loader/lib/loaders/stylePostLoader.js!./node_modules/postcss-loader/src/index.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/ProductVault.vue?vue&type=style&index=0&id=1a1232b4&lang=css":
/*!****************************************************************************************************************************************************************************************************************************************************************************************************************************!*\
  !*** ./node_modules/style-loader!./node_modules/css-loader??ref--5-1!./node_modules/vue-loader/lib/loaders/stylePostLoader.js!./node_modules/postcss-loader/src??ref--5-2!./node_modules/vue-loader/lib??vue-loader-options!./resources/js/components/salev2/ProductVault.vue?vue&type=style&index=0&id=1a1232b4&lang=css ***!
  \****************************************************************************************************************************************************************************************************************************************************************************************************************************/
/*! no static exports found */
/***/ (function(module, exports, __webpack_require__) {


var content = __webpack_require__(/*! !../../../../node_modules/css-loader??ref--5-1!../../../../node_modules/vue-loader/lib/loaders/stylePostLoader.js!../../../../node_modules/postcss-loader/src??ref--5-2!../../../../node_modules/vue-loader/lib??vue-loader-options!./ProductVault.vue?vue&type=style&index=0&id=1a1232b4&lang=css */ "./node_modules/css-loader/index.js?!./node_modules/vue-loader/lib/loaders/stylePostLoader.js!./node_modules/postcss-loader/src/index.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/ProductVault.vue?vue&type=style&index=0&id=1a1232b4&lang=css");

if(typeof content === 'string') content = [[module.i, content, '']];

var transform;
var insertInto;



var options = {"hmr":true}

options.transform = transform
options.insertInto = undefined;

var update = __webpack_require__(/*! ../../../../node_modules/style-loader/lib/addStyles.js */ "./node_modules/style-loader/lib/addStyles.js")(content, options);

if(content.locals) module.exports = content.locals;

if(false) {}

/***/ }),

/***/ "./node_modules/style-loader/index.js!./node_modules/css-loader/index.js?!./node_modules/vue-loader/lib/loaders/stylePostLoader.js!./node_modules/postcss-loader/src/index.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/SmartUpload.vue?vue&type=style&index=0&id=30d6d489&lang=css":
/*!***************************************************************************************************************************************************************************************************************************************************************************************************************************!*\
  !*** ./node_modules/style-loader!./node_modules/css-loader??ref--5-1!./node_modules/vue-loader/lib/loaders/stylePostLoader.js!./node_modules/postcss-loader/src??ref--5-2!./node_modules/vue-loader/lib??vue-loader-options!./resources/js/components/salev2/SmartUpload.vue?vue&type=style&index=0&id=30d6d489&lang=css ***!
  \***************************************************************************************************************************************************************************************************************************************************************************************************************************/
/*! no static exports found */
/***/ (function(module, exports, __webpack_require__) {


var content = __webpack_require__(/*! !../../../../node_modules/css-loader??ref--5-1!../../../../node_modules/vue-loader/lib/loaders/stylePostLoader.js!../../../../node_modules/postcss-loader/src??ref--5-2!../../../../node_modules/vue-loader/lib??vue-loader-options!./SmartUpload.vue?vue&type=style&index=0&id=30d6d489&lang=css */ "./node_modules/css-loader/index.js?!./node_modules/vue-loader/lib/loaders/stylePostLoader.js!./node_modules/postcss-loader/src/index.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/SmartUpload.vue?vue&type=style&index=0&id=30d6d489&lang=css");

if(typeof content === 'string') content = [[module.i, content, '']];

var transform;
var insertInto;



var options = {"hmr":true}

options.transform = transform
options.insertInto = undefined;

var update = __webpack_require__(/*! ../../../../node_modules/style-loader/lib/addStyles.js */ "./node_modules/style-loader/lib/addStyles.js")(content, options);

if(content.locals) module.exports = content.locals;

if(false) {}

/***/ }),

/***/ "./node_modules/style-loader/lib/addStyles.js":
/*!****************************************************!*\
  !*** ./node_modules/style-loader/lib/addStyles.js ***!
  \****************************************************/
/*! no static exports found */
/***/ (function(module, exports, __webpack_require__) {

/*
	MIT License http://www.opensource.org/licenses/mit-license.php
	Author Tobias Koppers @sokra
*/

var stylesInDom = {};

var	memoize = function (fn) {
	var memo;

	return function () {
		if (typeof memo === "undefined") memo = fn.apply(this, arguments);
		return memo;
	};
};

var isOldIE = memoize(function () {
	// Test for IE <= 9 as proposed by Browserhacks
	// @see http://browserhacks.com/#hack-e71d8692f65334173fee715c222cb805
	// Tests for existence of standard globals is to allow style-loader
	// to operate correctly into non-standard environments
	// @see https://github.com/webpack-contrib/style-loader/issues/177
	return window && document && document.all && !window.atob;
});

var getTarget = function (target, parent) {
  if (parent){
    return parent.querySelector(target);
  }
  return document.querySelector(target);
};

var getElement = (function (fn) {
	var memo = {};

	return function(target, parent) {
                // If passing function in options, then use it for resolve "head" element.
                // Useful for Shadow Root style i.e
                // {
                //   insertInto: function () { return document.querySelector("#foo").shadowRoot }
                // }
                if (typeof target === 'function') {
                        return target();
                }
                if (typeof memo[target] === "undefined") {
			var styleTarget = getTarget.call(this, target, parent);
			// Special case to return head of iframe instead of iframe itself
			if (window.HTMLIFrameElement && styleTarget instanceof window.HTMLIFrameElement) {
				try {
					// This will throw an exception if access to iframe is blocked
					// due to cross-origin restrictions
					styleTarget = styleTarget.contentDocument.head;
				} catch(e) {
					styleTarget = null;
				}
			}
			memo[target] = styleTarget;
		}
		return memo[target]
	};
})();

var singleton = null;
var	singletonCounter = 0;
var	stylesInsertedAtTop = [];

var	fixUrls = __webpack_require__(/*! ./urls */ "./node_modules/style-loader/lib/urls.js");

module.exports = function(list, options) {
	if (typeof DEBUG !== "undefined" && DEBUG) {
		if (typeof document !== "object") throw new Error("The style-loader cannot be used in a non-browser environment");
	}

	options = options || {};

	options.attrs = typeof options.attrs === "object" ? options.attrs : {};

	// Force single-tag solution on IE6-9, which has a hard limit on the # of <style>
	// tags it will allow on a page
	if (!options.singleton && typeof options.singleton !== "boolean") options.singleton = isOldIE();

	// By default, add <style> tags to the <head> element
        if (!options.insertInto) options.insertInto = "head";

	// By default, add <style> tags to the bottom of the target
	if (!options.insertAt) options.insertAt = "bottom";

	var styles = listToStyles(list, options);

	addStylesToDom(styles, options);

	return function update (newList) {
		var mayRemove = [];

		for (var i = 0; i < styles.length; i++) {
			var item = styles[i];
			var domStyle = stylesInDom[item.id];

			domStyle.refs--;
			mayRemove.push(domStyle);
		}

		if(newList) {
			var newStyles = listToStyles(newList, options);
			addStylesToDom(newStyles, options);
		}

		for (var i = 0; i < mayRemove.length; i++) {
			var domStyle = mayRemove[i];

			if(domStyle.refs === 0) {
				for (var j = 0; j < domStyle.parts.length; j++) domStyle.parts[j]();

				delete stylesInDom[domStyle.id];
			}
		}
	};
};

function addStylesToDom (styles, options) {
	for (var i = 0; i < styles.length; i++) {
		var item = styles[i];
		var domStyle = stylesInDom[item.id];

		if(domStyle) {
			domStyle.refs++;

			for(var j = 0; j < domStyle.parts.length; j++) {
				domStyle.parts[j](item.parts[j]);
			}

			for(; j < item.parts.length; j++) {
				domStyle.parts.push(addStyle(item.parts[j], options));
			}
		} else {
			var parts = [];

			for(var j = 0; j < item.parts.length; j++) {
				parts.push(addStyle(item.parts[j], options));
			}

			stylesInDom[item.id] = {id: item.id, refs: 1, parts: parts};
		}
	}
}

function listToStyles (list, options) {
	var styles = [];
	var newStyles = {};

	for (var i = 0; i < list.length; i++) {
		var item = list[i];
		var id = options.base ? item[0] + options.base : item[0];
		var css = item[1];
		var media = item[2];
		var sourceMap = item[3];
		var part = {css: css, media: media, sourceMap: sourceMap};

		if(!newStyles[id]) styles.push(newStyles[id] = {id: id, parts: [part]});
		else newStyles[id].parts.push(part);
	}

	return styles;
}

function insertStyleElement (options, style) {
	var target = getElement(options.insertInto)

	if (!target) {
		throw new Error("Couldn't find a style target. This probably means that the value for the 'insertInto' parameter is invalid.");
	}

	var lastStyleElementInsertedAtTop = stylesInsertedAtTop[stylesInsertedAtTop.length - 1];

	if (options.insertAt === "top") {
		if (!lastStyleElementInsertedAtTop) {
			target.insertBefore(style, target.firstChild);
		} else if (lastStyleElementInsertedAtTop.nextSibling) {
			target.insertBefore(style, lastStyleElementInsertedAtTop.nextSibling);
		} else {
			target.appendChild(style);
		}
		stylesInsertedAtTop.push(style);
	} else if (options.insertAt === "bottom") {
		target.appendChild(style);
	} else if (typeof options.insertAt === "object" && options.insertAt.before) {
		var nextSibling = getElement(options.insertAt.before, target);
		target.insertBefore(style, nextSibling);
	} else {
		throw new Error("[Style Loader]\n\n Invalid value for parameter 'insertAt' ('options.insertAt') found.\n Must be 'top', 'bottom', or Object.\n (https://github.com/webpack-contrib/style-loader#insertat)\n");
	}
}

function removeStyleElement (style) {
	if (style.parentNode === null) return false;
	style.parentNode.removeChild(style);

	var idx = stylesInsertedAtTop.indexOf(style);
	if(idx >= 0) {
		stylesInsertedAtTop.splice(idx, 1);
	}
}

function createStyleElement (options) {
	var style = document.createElement("style");

	if(options.attrs.type === undefined) {
		options.attrs.type = "text/css";
	}

	if(options.attrs.nonce === undefined) {
		var nonce = getNonce();
		if (nonce) {
			options.attrs.nonce = nonce;
		}
	}

	addAttrs(style, options.attrs);
	insertStyleElement(options, style);

	return style;
}

function createLinkElement (options) {
	var link = document.createElement("link");

	if(options.attrs.type === undefined) {
		options.attrs.type = "text/css";
	}
	options.attrs.rel = "stylesheet";

	addAttrs(link, options.attrs);
	insertStyleElement(options, link);

	return link;
}

function addAttrs (el, attrs) {
	Object.keys(attrs).forEach(function (key) {
		el.setAttribute(key, attrs[key]);
	});
}

function getNonce() {
	if (false) {}

	return __webpack_require__.nc;
}

function addStyle (obj, options) {
	var style, update, remove, result;

	// If a transform function was defined, run it on the css
	if (options.transform && obj.css) {
	    result = typeof options.transform === 'function'
		 ? options.transform(obj.css) 
		 : options.transform.default(obj.css);

	    if (result) {
	    	// If transform returns a value, use that instead of the original css.
	    	// This allows running runtime transformations on the css.
	    	obj.css = result;
	    } else {
	    	// If the transform function returns a falsy value, don't add this css.
	    	// This allows conditional loading of css
	    	return function() {
	    		// noop
	    	};
	    }
	}

	if (options.singleton) {
		var styleIndex = singletonCounter++;

		style = singleton || (singleton = createStyleElement(options));

		update = applyToSingletonTag.bind(null, style, styleIndex, false);
		remove = applyToSingletonTag.bind(null, style, styleIndex, true);

	} else if (
		obj.sourceMap &&
		typeof URL === "function" &&
		typeof URL.createObjectURL === "function" &&
		typeof URL.revokeObjectURL === "function" &&
		typeof Blob === "function" &&
		typeof btoa === "function"
	) {
		style = createLinkElement(options);
		update = updateLink.bind(null, style, options);
		remove = function () {
			removeStyleElement(style);

			if(style.href) URL.revokeObjectURL(style.href);
		};
	} else {
		style = createStyleElement(options);
		update = applyToTag.bind(null, style);
		remove = function () {
			removeStyleElement(style);
		};
	}

	update(obj);

	return function updateStyle (newObj) {
		if (newObj) {
			if (
				newObj.css === obj.css &&
				newObj.media === obj.media &&
				newObj.sourceMap === obj.sourceMap
			) {
				return;
			}

			update(obj = newObj);
		} else {
			remove();
		}
	};
}

var replaceText = (function () {
	var textStore = [];

	return function (index, replacement) {
		textStore[index] = replacement;

		return textStore.filter(Boolean).join('\n');
	};
})();

function applyToSingletonTag (style, index, remove, obj) {
	var css = remove ? "" : obj.css;

	if (style.styleSheet) {
		style.styleSheet.cssText = replaceText(index, css);
	} else {
		var cssNode = document.createTextNode(css);
		var childNodes = style.childNodes;

		if (childNodes[index]) style.removeChild(childNodes[index]);

		if (childNodes.length) {
			style.insertBefore(cssNode, childNodes[index]);
		} else {
			style.appendChild(cssNode);
		}
	}
}

function applyToTag (style, obj) {
	var css = obj.css;
	var media = obj.media;

	if(media) {
		style.setAttribute("media", media)
	}

	if(style.styleSheet) {
		style.styleSheet.cssText = css;
	} else {
		while(style.firstChild) {
			style.removeChild(style.firstChild);
		}

		style.appendChild(document.createTextNode(css));
	}
}

function updateLink (link, options, obj) {
	var css = obj.css;
	var sourceMap = obj.sourceMap;

	/*
		If convertToAbsoluteUrls isn't defined, but sourcemaps are enabled
		and there is no publicPath defined then lets turn convertToAbsoluteUrls
		on by default.  Otherwise default to the convertToAbsoluteUrls option
		directly
	*/
	var autoFixUrls = options.convertToAbsoluteUrls === undefined && sourceMap;

	if (options.convertToAbsoluteUrls || autoFixUrls) {
		css = fixUrls(css);
	}

	if (sourceMap) {
		// http://stackoverflow.com/a/26603875
		css += "\n/*# sourceMappingURL=data:application/json;base64," + btoa(unescape(encodeURIComponent(JSON.stringify(sourceMap)))) + " */";
	}

	var blob = new Blob([css], { type: "text/css" });

	var oldSrc = link.href;

	link.href = URL.createObjectURL(blob);

	if(oldSrc) URL.revokeObjectURL(oldSrc);
}


/***/ }),

/***/ "./node_modules/style-loader/lib/urls.js":
/*!***********************************************!*\
  !*** ./node_modules/style-loader/lib/urls.js ***!
  \***********************************************/
/*! no static exports found */
/***/ (function(module, exports) {


/**
 * When source maps are enabled, `style-loader` uses a link element with a data-uri to
 * embed the css on the page. This breaks all relative urls because now they are relative to a
 * bundle instead of the current page.
 *
 * One solution is to only use full urls, but that may be impossible.
 *
 * Instead, this function "fixes" the relative urls to be absolute according to the current page location.
 *
 * A rudimentary test suite is located at `test/fixUrls.js` and can be run via the `npm test` command.
 *
 */

module.exports = function (css) {
  // get current location
  var location = typeof window !== "undefined" && window.location;

  if (!location) {
    throw new Error("fixUrls requires window.location");
  }

	// blank or null?
	if (!css || typeof css !== "string") {
	  return css;
  }

  var baseUrl = location.protocol + "//" + location.host;
  var currentDir = baseUrl + location.pathname.replace(/\/[^\/]*$/, "/");

	// convert each url(...)
	/*
	This regular expression is just a way to recursively match brackets within
	a string.

	 /url\s*\(  = Match on the word "url" with any whitespace after it and then a parens
	   (  = Start a capturing group
	     (?:  = Start a non-capturing group
	         [^)(]  = Match anything that isn't a parentheses
	         |  = OR
	         \(  = Match a start parentheses
	             (?:  = Start another non-capturing groups
	                 [^)(]+  = Match anything that isn't a parentheses
	                 |  = OR
	                 \(  = Match a start parentheses
	                     [^)(]*  = Match anything that isn't a parentheses
	                 \)  = Match a end parentheses
	             )  = End Group
              *\) = Match anything and then a close parens
          )  = Close non-capturing group
          *  = Match anything
       )  = Close capturing group
	 \)  = Match a close parens

	 /gi  = Get all matches, not the first.  Be case insensitive.
	 */
	var fixedCss = css.replace(/url\s*\(((?:[^)(]|\((?:[^)(]+|\([^)(]*\))*\))*)\)/gi, function(fullMatch, origUrl) {
		// strip quotes (if they exist)
		var unquotedOrigUrl = origUrl
			.trim()
			.replace(/^"(.*)"$/, function(o, $1){ return $1; })
			.replace(/^'(.*)'$/, function(o, $1){ return $1; });

		// already a full url? no change
		if (/^(#|data:|http:\/\/|https:\/\/|file:\/\/\/|\s*$)/i.test(unquotedOrigUrl)) {
		  return fullMatch;
		}

		// convert the url to a full url
		var newUrl;

		if (unquotedOrigUrl.indexOf("//") === 0) {
		  	//TODO: should we add protocol?
			newUrl = unquotedOrigUrl;
		} else if (unquotedOrigUrl.indexOf("/") === 0) {
			// path should be relative to the base url
			newUrl = baseUrl + unquotedOrigUrl; // already starts with '/'
		} else {
			// path should be relative to current directory
			newUrl = currentDir + unquotedOrigUrl.replace(/^\.\//, ""); // Strip leading './'
		}

		// send back the fixed url(...)
		return "url(" + JSON.stringify(newUrl) + ")";
	});

	// send back the fixed css
	return fixedCss;
};


/***/ }),

/***/ "./node_modules/vue-loader/lib/runtime/componentNormalizer.js":
/*!********************************************************************!*\
  !*** ./node_modules/vue-loader/lib/runtime/componentNormalizer.js ***!
  \********************************************************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "default", function() { return normalizeComponent; });
/* globals __VUE_SSR_CONTEXT__ */

// IMPORTANT: Do NOT use ES2015 features in this file (except for modules).
// This module is a runtime utility for cleaner component module output and will
// be included in the final webpack user bundle.

function normalizeComponent(
  scriptExports,
  render,
  staticRenderFns,
  functionalTemplate,
  injectStyles,
  scopeId,
  moduleIdentifier /* server only */,
  shadowMode /* vue-cli only */
) {
  // Vue.extend constructor export interop
  var options =
    typeof scriptExports === 'function' ? scriptExports.options : scriptExports

  // render functions
  if (render) {
    options.render = render
    options.staticRenderFns = staticRenderFns
    options._compiled = true
  }

  // functional template
  if (functionalTemplate) {
    options.functional = true
  }

  // scopedId
  if (scopeId) {
    options._scopeId = 'data-v-' + scopeId
  }

  var hook
  if (moduleIdentifier) {
    // server build
    hook = function (context) {
      // 2.3 injection
      context =
        context || // cached call
        (this.$vnode && this.$vnode.ssrContext) || // stateful
        (this.parent && this.parent.$vnode && this.parent.$vnode.ssrContext) // functional
      // 2.2 with runInNewContext: true
      if (!context && typeof __VUE_SSR_CONTEXT__ !== 'undefined') {
        context = __VUE_SSR_CONTEXT__
      }
      // inject component styles
      if (injectStyles) {
        injectStyles.call(this, context)
      }
      // register component module identifier for async chunk inferrence
      if (context && context._registeredComponents) {
        context._registeredComponents.add(moduleIdentifier)
      }
    }
    // used by ssr in case component is cached and beforeCreate
    // never gets called
    options._ssrRegister = hook
  } else if (injectStyles) {
    hook = shadowMode
      ? function () {
          injectStyles.call(
            this,
            (options.functional ? this.parent : this).$root.$options.shadowRoot
          )
        }
      : injectStyles
  }

  if (hook) {
    if (options.functional) {
      // for template-only hot-reload because in that case the render fn doesn't
      // go through the normalizer
      options._injectStyles = hook
      // register for functional component in vue file
      var originalRender = options.render
      options.render = function renderWithStyleInjection(h, context) {
        hook.call(context)
        return originalRender(h, context)
      }
    } else {
      // inject component registration as beforeCreate hook
      var existing = options.beforeCreate
      options.beforeCreate = existing ? [].concat(existing, hook) : [hook]
    }
  }

  return {
    exports: scriptExports,
    options: options
  }
}


/***/ }),

/***/ "./resources/js/app.js":
/*!*****************************!*\
  !*** ./resources/js/app.js ***!
  \*****************************/
/*! no exports provided */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var _components_salev2_Dashboard_vue__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ./components/salev2/Dashboard.vue */ "./resources/js/components/salev2/Dashboard.vue");
/* harmony import */ var _components_salev2_ProductVault_vue__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./components/salev2/ProductVault.vue */ "./resources/js/components/salev2/ProductVault.vue");
/* harmony import */ var _components_salev2_OrderManager_vue__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ./components/salev2/OrderManager.vue */ "./resources/js/components/salev2/OrderManager.vue");
/* harmony import */ var _components_salev2_SmartUpload_vue__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! ./components/salev2/SmartUpload.vue */ "./resources/js/components/salev2/SmartUpload.vue");
/* harmony import */ var _components_salev2_CustomerManager_vue__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(/*! ./components/salev2/CustomerManager.vue */ "./resources/js/components/salev2/CustomerManager.vue");
/* harmony import */ var _components_salev2_OrderCreate_vue__WEBPACK_IMPORTED_MODULE_5__ = __webpack_require__(/*! ./components/salev2/OrderCreate.vue */ "./resources/js/components/salev2/OrderCreate.vue");
/* harmony import */ var _components_salev2_ProductCreate_vue__WEBPACK_IMPORTED_MODULE_6__ = __webpack_require__(/*! ./components/salev2/ProductCreate.vue */ "./resources/js/components/salev2/ProductCreate.vue");
/* harmony import */ var _components_salev2_ProductEdit_vue__WEBPACK_IMPORTED_MODULE_7__ = __webpack_require__(/*! ./components/salev2/ProductEdit.vue */ "./resources/js/components/salev2/ProductEdit.vue");
// Map global CDN variables if they use different names
var ElementUI = window.ELEMENT;
var VueRouter = window.VueRouter;
var VueApexCharts = window.VueApexCharts;
var axios = window.axios;
if (typeof Vue !== 'undefined') {
  if (typeof VueRouter !== 'undefined') Vue.use(VueRouter);
  if (typeof ElementUI !== 'undefined') Vue.use(ElementUI);
  if (typeof VueApexCharts !== 'undefined') Vue.component('apexchart', VueApexCharts);
}

// Components








var routes = [{
  path: '/salev2',
  component: _components_salev2_Dashboard_vue__WEBPACK_IMPORTED_MODULE_0__["default"]
}, {
  path: '/salev2/products',
  component: _components_salev2_ProductVault_vue__WEBPACK_IMPORTED_MODULE_1__["default"]
}, {
  path: '/salev2/orders',
  component: _components_salev2_OrderManager_vue__WEBPACK_IMPORTED_MODULE_2__["default"]
}, {
  path: '/salev2/upload',
  component: _components_salev2_SmartUpload_vue__WEBPACK_IMPORTED_MODULE_3__["default"]
}, {
  path: '/salev2/customers',
  component: _components_salev2_CustomerManager_vue__WEBPACK_IMPORTED_MODULE_4__["default"]
}, {
  path: '/salev2/order/create',
  name: 'order-create',
  component: _components_salev2_OrderCreate_vue__WEBPACK_IMPORTED_MODULE_5__["default"]
}, {
  path: '/salev2/order/:id/edit',
  name: 'order-edit',
  component: __webpack_require__(/*! ./components/salev2/OrderEdit.vue */ "./resources/js/components/salev2/OrderEdit.vue")["default"]
}, {
  path: '/salev2/product/create',
  name: 'product-create',
  component: _components_salev2_ProductCreate_vue__WEBPACK_IMPORTED_MODULE_6__["default"]
}, {
  path: '/salev2/product/:id/edit',
  name: 'product-edit',
  component: _components_salev2_ProductEdit_vue__WEBPACK_IMPORTED_MODULE_7__["default"]
}, {
  path: '*',
  redirect: '/salev2'
}];
var router = new VueRouter({
  mode: 'history',
  routes: routes
});

// CSRF for axios
axios.defaults.headers.common['X-Requested-With'] = 'XMLHttpRequest';
var token = document.head.querySelector('meta[name="csrf-token"]');
if (token) {
  axios.defaults.headers.common['X-CSRF-TOKEN'] = token.content;
}
var app = new Vue({
  el: '#app',
  router: router
});

/***/ }),

/***/ "./resources/js/components/salev2/CustomerManager.vue":
/*!************************************************************!*\
  !*** ./resources/js/components/salev2/CustomerManager.vue ***!
  \************************************************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var _CustomerManager_vue_vue_type_template_id_1e0b8b4e__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ./CustomerManager.vue?vue&type=template&id=1e0b8b4e */ "./resources/js/components/salev2/CustomerManager.vue?vue&type=template&id=1e0b8b4e");
/* harmony import */ var _CustomerManager_vue_vue_type_script_lang_js__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./CustomerManager.vue?vue&type=script&lang=js */ "./resources/js/components/salev2/CustomerManager.vue?vue&type=script&lang=js");
/* empty/unused harmony star reexport *//* harmony import */ var _node_modules_vue_loader_lib_runtime_componentNormalizer_js__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ../../../../node_modules/vue-loader/lib/runtime/componentNormalizer.js */ "./node_modules/vue-loader/lib/runtime/componentNormalizer.js");





/* normalize component */

var component = Object(_node_modules_vue_loader_lib_runtime_componentNormalizer_js__WEBPACK_IMPORTED_MODULE_2__["default"])(
  _CustomerManager_vue_vue_type_script_lang_js__WEBPACK_IMPORTED_MODULE_1__["default"],
  _CustomerManager_vue_vue_type_template_id_1e0b8b4e__WEBPACK_IMPORTED_MODULE_0__["render"],
  _CustomerManager_vue_vue_type_template_id_1e0b8b4e__WEBPACK_IMPORTED_MODULE_0__["staticRenderFns"],
  false,
  null,
  null,
  null
  
)

/* hot reload */
if (false) { var api; }
component.options.__file = "resources/js/components/salev2/CustomerManager.vue"
/* harmony default export */ __webpack_exports__["default"] = (component.exports);

/***/ }),

/***/ "./resources/js/components/salev2/CustomerManager.vue?vue&type=script&lang=js":
/*!************************************************************************************!*\
  !*** ./resources/js/components/salev2/CustomerManager.vue?vue&type=script&lang=js ***!
  \************************************************************************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var _node_modules_babel_loader_lib_index_js_ref_4_0_node_modules_vue_loader_lib_index_js_vue_loader_options_CustomerManager_vue_vue_type_script_lang_js__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! -!../../../../node_modules/babel-loader/lib??ref--4-0!../../../../node_modules/vue-loader/lib??vue-loader-options!./CustomerManager.vue?vue&type=script&lang=js */ "./node_modules/babel-loader/lib/index.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/CustomerManager.vue?vue&type=script&lang=js");
/* empty/unused harmony star reexport */ /* harmony default export */ __webpack_exports__["default"] = (_node_modules_babel_loader_lib_index_js_ref_4_0_node_modules_vue_loader_lib_index_js_vue_loader_options_CustomerManager_vue_vue_type_script_lang_js__WEBPACK_IMPORTED_MODULE_0__["default"]); 

/***/ }),

/***/ "./resources/js/components/salev2/CustomerManager.vue?vue&type=template&id=1e0b8b4e":
/*!******************************************************************************************!*\
  !*** ./resources/js/components/salev2/CustomerManager.vue?vue&type=template&id=1e0b8b4e ***!
  \******************************************************************************************/
/*! exports provided: render, staticRenderFns */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var _node_modules_babel_loader_lib_index_js_ref_4_0_node_modules_vue_loader_lib_loaders_templateLoader_js_ref_6_node_modules_vue_loader_lib_index_js_vue_loader_options_CustomerManager_vue_vue_type_template_id_1e0b8b4e__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! -!../../../../node_modules/babel-loader/lib??ref--4-0!../../../../node_modules/vue-loader/lib/loaders/templateLoader.js??ref--6!../../../../node_modules/vue-loader/lib??vue-loader-options!./CustomerManager.vue?vue&type=template&id=1e0b8b4e */ "./node_modules/babel-loader/lib/index.js?!./node_modules/vue-loader/lib/loaders/templateLoader.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/CustomerManager.vue?vue&type=template&id=1e0b8b4e");
/* harmony reexport (safe) */ __webpack_require__.d(__webpack_exports__, "render", function() { return _node_modules_babel_loader_lib_index_js_ref_4_0_node_modules_vue_loader_lib_loaders_templateLoader_js_ref_6_node_modules_vue_loader_lib_index_js_vue_loader_options_CustomerManager_vue_vue_type_template_id_1e0b8b4e__WEBPACK_IMPORTED_MODULE_0__["render"]; });

/* harmony reexport (safe) */ __webpack_require__.d(__webpack_exports__, "staticRenderFns", function() { return _node_modules_babel_loader_lib_index_js_ref_4_0_node_modules_vue_loader_lib_loaders_templateLoader_js_ref_6_node_modules_vue_loader_lib_index_js_vue_loader_options_CustomerManager_vue_vue_type_template_id_1e0b8b4e__WEBPACK_IMPORTED_MODULE_0__["staticRenderFns"]; });



/***/ }),

/***/ "./resources/js/components/salev2/Dashboard.vue":
/*!******************************************************!*\
  !*** ./resources/js/components/salev2/Dashboard.vue ***!
  \******************************************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var _Dashboard_vue_vue_type_template_id_16392a9a__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ./Dashboard.vue?vue&type=template&id=16392a9a */ "./resources/js/components/salev2/Dashboard.vue?vue&type=template&id=16392a9a");
/* harmony import */ var _Dashboard_vue_vue_type_script_lang_js__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./Dashboard.vue?vue&type=script&lang=js */ "./resources/js/components/salev2/Dashboard.vue?vue&type=script&lang=js");
/* empty/unused harmony star reexport *//* harmony import */ var _Dashboard_vue_vue_type_style_index_0_id_16392a9a_lang_css__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ./Dashboard.vue?vue&type=style&index=0&id=16392a9a&lang=css */ "./resources/js/components/salev2/Dashboard.vue?vue&type=style&index=0&id=16392a9a&lang=css");
/* harmony import */ var _node_modules_vue_loader_lib_runtime_componentNormalizer_js__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! ../../../../node_modules/vue-loader/lib/runtime/componentNormalizer.js */ "./node_modules/vue-loader/lib/runtime/componentNormalizer.js");






/* normalize component */

var component = Object(_node_modules_vue_loader_lib_runtime_componentNormalizer_js__WEBPACK_IMPORTED_MODULE_3__["default"])(
  _Dashboard_vue_vue_type_script_lang_js__WEBPACK_IMPORTED_MODULE_1__["default"],
  _Dashboard_vue_vue_type_template_id_16392a9a__WEBPACK_IMPORTED_MODULE_0__["render"],
  _Dashboard_vue_vue_type_template_id_16392a9a__WEBPACK_IMPORTED_MODULE_0__["staticRenderFns"],
  false,
  null,
  null,
  null
  
)

/* hot reload */
if (false) { var api; }
component.options.__file = "resources/js/components/salev2/Dashboard.vue"
/* harmony default export */ __webpack_exports__["default"] = (component.exports);

/***/ }),

/***/ "./resources/js/components/salev2/Dashboard.vue?vue&type=script&lang=js":
/*!******************************************************************************!*\
  !*** ./resources/js/components/salev2/Dashboard.vue?vue&type=script&lang=js ***!
  \******************************************************************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var _node_modules_babel_loader_lib_index_js_ref_4_0_node_modules_vue_loader_lib_index_js_vue_loader_options_Dashboard_vue_vue_type_script_lang_js__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! -!../../../../node_modules/babel-loader/lib??ref--4-0!../../../../node_modules/vue-loader/lib??vue-loader-options!./Dashboard.vue?vue&type=script&lang=js */ "./node_modules/babel-loader/lib/index.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/Dashboard.vue?vue&type=script&lang=js");
/* empty/unused harmony star reexport */ /* harmony default export */ __webpack_exports__["default"] = (_node_modules_babel_loader_lib_index_js_ref_4_0_node_modules_vue_loader_lib_index_js_vue_loader_options_Dashboard_vue_vue_type_script_lang_js__WEBPACK_IMPORTED_MODULE_0__["default"]); 

/***/ }),

/***/ "./resources/js/components/salev2/Dashboard.vue?vue&type=style&index=0&id=16392a9a&lang=css":
/*!**************************************************************************************************!*\
  !*** ./resources/js/components/salev2/Dashboard.vue?vue&type=style&index=0&id=16392a9a&lang=css ***!
  \**************************************************************************************************/
/*! no static exports found */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var _node_modules_style_loader_index_js_node_modules_css_loader_index_js_ref_5_1_node_modules_vue_loader_lib_loaders_stylePostLoader_js_node_modules_postcss_loader_src_index_js_ref_5_2_node_modules_vue_loader_lib_index_js_vue_loader_options_Dashboard_vue_vue_type_style_index_0_id_16392a9a_lang_css__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! -!../../../../node_modules/style-loader!../../../../node_modules/css-loader??ref--5-1!../../../../node_modules/vue-loader/lib/loaders/stylePostLoader.js!../../../../node_modules/postcss-loader/src??ref--5-2!../../../../node_modules/vue-loader/lib??vue-loader-options!./Dashboard.vue?vue&type=style&index=0&id=16392a9a&lang=css */ "./node_modules/style-loader/index.js!./node_modules/css-loader/index.js?!./node_modules/vue-loader/lib/loaders/stylePostLoader.js!./node_modules/postcss-loader/src/index.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/Dashboard.vue?vue&type=style&index=0&id=16392a9a&lang=css");
/* harmony import */ var _node_modules_style_loader_index_js_node_modules_css_loader_index_js_ref_5_1_node_modules_vue_loader_lib_loaders_stylePostLoader_js_node_modules_postcss_loader_src_index_js_ref_5_2_node_modules_vue_loader_lib_index_js_vue_loader_options_Dashboard_vue_vue_type_style_index_0_id_16392a9a_lang_css__WEBPACK_IMPORTED_MODULE_0___default = /*#__PURE__*/__webpack_require__.n(_node_modules_style_loader_index_js_node_modules_css_loader_index_js_ref_5_1_node_modules_vue_loader_lib_loaders_stylePostLoader_js_node_modules_postcss_loader_src_index_js_ref_5_2_node_modules_vue_loader_lib_index_js_vue_loader_options_Dashboard_vue_vue_type_style_index_0_id_16392a9a_lang_css__WEBPACK_IMPORTED_MODULE_0__);
/* harmony reexport (unknown) */ for(var __WEBPACK_IMPORT_KEY__ in _node_modules_style_loader_index_js_node_modules_css_loader_index_js_ref_5_1_node_modules_vue_loader_lib_loaders_stylePostLoader_js_node_modules_postcss_loader_src_index_js_ref_5_2_node_modules_vue_loader_lib_index_js_vue_loader_options_Dashboard_vue_vue_type_style_index_0_id_16392a9a_lang_css__WEBPACK_IMPORTED_MODULE_0__) if(["default"].indexOf(__WEBPACK_IMPORT_KEY__) < 0) (function(key) { __webpack_require__.d(__webpack_exports__, key, function() { return _node_modules_style_loader_index_js_node_modules_css_loader_index_js_ref_5_1_node_modules_vue_loader_lib_loaders_stylePostLoader_js_node_modules_postcss_loader_src_index_js_ref_5_2_node_modules_vue_loader_lib_index_js_vue_loader_options_Dashboard_vue_vue_type_style_index_0_id_16392a9a_lang_css__WEBPACK_IMPORTED_MODULE_0__[key]; }) }(__WEBPACK_IMPORT_KEY__));


/***/ }),

/***/ "./resources/js/components/salev2/Dashboard.vue?vue&type=template&id=16392a9a":
/*!************************************************************************************!*\
  !*** ./resources/js/components/salev2/Dashboard.vue?vue&type=template&id=16392a9a ***!
  \************************************************************************************/
/*! exports provided: render, staticRenderFns */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var _node_modules_babel_loader_lib_index_js_ref_4_0_node_modules_vue_loader_lib_loaders_templateLoader_js_ref_6_node_modules_vue_loader_lib_index_js_vue_loader_options_Dashboard_vue_vue_type_template_id_16392a9a__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! -!../../../../node_modules/babel-loader/lib??ref--4-0!../../../../node_modules/vue-loader/lib/loaders/templateLoader.js??ref--6!../../../../node_modules/vue-loader/lib??vue-loader-options!./Dashboard.vue?vue&type=template&id=16392a9a */ "./node_modules/babel-loader/lib/index.js?!./node_modules/vue-loader/lib/loaders/templateLoader.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/Dashboard.vue?vue&type=template&id=16392a9a");
/* harmony reexport (safe) */ __webpack_require__.d(__webpack_exports__, "render", function() { return _node_modules_babel_loader_lib_index_js_ref_4_0_node_modules_vue_loader_lib_loaders_templateLoader_js_ref_6_node_modules_vue_loader_lib_index_js_vue_loader_options_Dashboard_vue_vue_type_template_id_16392a9a__WEBPACK_IMPORTED_MODULE_0__["render"]; });

/* harmony reexport (safe) */ __webpack_require__.d(__webpack_exports__, "staticRenderFns", function() { return _node_modules_babel_loader_lib_index_js_ref_4_0_node_modules_vue_loader_lib_loaders_templateLoader_js_ref_6_node_modules_vue_loader_lib_index_js_vue_loader_options_Dashboard_vue_vue_type_template_id_16392a9a__WEBPACK_IMPORTED_MODULE_0__["staticRenderFns"]; });



/***/ }),

/***/ "./resources/js/components/salev2/OrderCreate.vue":
/*!********************************************************!*\
  !*** ./resources/js/components/salev2/OrderCreate.vue ***!
  \********************************************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var _OrderCreate_vue_vue_type_template_id_5edbfd6e__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ./OrderCreate.vue?vue&type=template&id=5edbfd6e */ "./resources/js/components/salev2/OrderCreate.vue?vue&type=template&id=5edbfd6e");
/* harmony import */ var _OrderCreate_vue_vue_type_script_lang_js__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./OrderCreate.vue?vue&type=script&lang=js */ "./resources/js/components/salev2/OrderCreate.vue?vue&type=script&lang=js");
/* empty/unused harmony star reexport *//* harmony import */ var _node_modules_vue_loader_lib_runtime_componentNormalizer_js__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ../../../../node_modules/vue-loader/lib/runtime/componentNormalizer.js */ "./node_modules/vue-loader/lib/runtime/componentNormalizer.js");





/* normalize component */

var component = Object(_node_modules_vue_loader_lib_runtime_componentNormalizer_js__WEBPACK_IMPORTED_MODULE_2__["default"])(
  _OrderCreate_vue_vue_type_script_lang_js__WEBPACK_IMPORTED_MODULE_1__["default"],
  _OrderCreate_vue_vue_type_template_id_5edbfd6e__WEBPACK_IMPORTED_MODULE_0__["render"],
  _OrderCreate_vue_vue_type_template_id_5edbfd6e__WEBPACK_IMPORTED_MODULE_0__["staticRenderFns"],
  false,
  null,
  null,
  null
  
)

/* hot reload */
if (false) { var api; }
component.options.__file = "resources/js/components/salev2/OrderCreate.vue"
/* harmony default export */ __webpack_exports__["default"] = (component.exports);

/***/ }),

/***/ "./resources/js/components/salev2/OrderCreate.vue?vue&type=script&lang=js":
/*!********************************************************************************!*\
  !*** ./resources/js/components/salev2/OrderCreate.vue?vue&type=script&lang=js ***!
  \********************************************************************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var _node_modules_babel_loader_lib_index_js_ref_4_0_node_modules_vue_loader_lib_index_js_vue_loader_options_OrderCreate_vue_vue_type_script_lang_js__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! -!../../../../node_modules/babel-loader/lib??ref--4-0!../../../../node_modules/vue-loader/lib??vue-loader-options!./OrderCreate.vue?vue&type=script&lang=js */ "./node_modules/babel-loader/lib/index.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/OrderCreate.vue?vue&type=script&lang=js");
/* empty/unused harmony star reexport */ /* harmony default export */ __webpack_exports__["default"] = (_node_modules_babel_loader_lib_index_js_ref_4_0_node_modules_vue_loader_lib_index_js_vue_loader_options_OrderCreate_vue_vue_type_script_lang_js__WEBPACK_IMPORTED_MODULE_0__["default"]); 

/***/ }),

/***/ "./resources/js/components/salev2/OrderCreate.vue?vue&type=template&id=5edbfd6e":
/*!**************************************************************************************!*\
  !*** ./resources/js/components/salev2/OrderCreate.vue?vue&type=template&id=5edbfd6e ***!
  \**************************************************************************************/
/*! exports provided: render, staticRenderFns */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var _node_modules_babel_loader_lib_index_js_ref_4_0_node_modules_vue_loader_lib_loaders_templateLoader_js_ref_6_node_modules_vue_loader_lib_index_js_vue_loader_options_OrderCreate_vue_vue_type_template_id_5edbfd6e__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! -!../../../../node_modules/babel-loader/lib??ref--4-0!../../../../node_modules/vue-loader/lib/loaders/templateLoader.js??ref--6!../../../../node_modules/vue-loader/lib??vue-loader-options!./OrderCreate.vue?vue&type=template&id=5edbfd6e */ "./node_modules/babel-loader/lib/index.js?!./node_modules/vue-loader/lib/loaders/templateLoader.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/OrderCreate.vue?vue&type=template&id=5edbfd6e");
/* harmony reexport (safe) */ __webpack_require__.d(__webpack_exports__, "render", function() { return _node_modules_babel_loader_lib_index_js_ref_4_0_node_modules_vue_loader_lib_loaders_templateLoader_js_ref_6_node_modules_vue_loader_lib_index_js_vue_loader_options_OrderCreate_vue_vue_type_template_id_5edbfd6e__WEBPACK_IMPORTED_MODULE_0__["render"]; });

/* harmony reexport (safe) */ __webpack_require__.d(__webpack_exports__, "staticRenderFns", function() { return _node_modules_babel_loader_lib_index_js_ref_4_0_node_modules_vue_loader_lib_loaders_templateLoader_js_ref_6_node_modules_vue_loader_lib_index_js_vue_loader_options_OrderCreate_vue_vue_type_template_id_5edbfd6e__WEBPACK_IMPORTED_MODULE_0__["staticRenderFns"]; });



/***/ }),

/***/ "./resources/js/components/salev2/OrderEdit.vue":
/*!******************************************************!*\
  !*** ./resources/js/components/salev2/OrderEdit.vue ***!
  \******************************************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var _OrderEdit_vue_vue_type_template_id_2ac7fed2__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ./OrderEdit.vue?vue&type=template&id=2ac7fed2 */ "./resources/js/components/salev2/OrderEdit.vue?vue&type=template&id=2ac7fed2");
/* harmony import */ var _OrderEdit_vue_vue_type_script_lang_js__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./OrderEdit.vue?vue&type=script&lang=js */ "./resources/js/components/salev2/OrderEdit.vue?vue&type=script&lang=js");
/* empty/unused harmony star reexport *//* harmony import */ var _node_modules_vue_loader_lib_runtime_componentNormalizer_js__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ../../../../node_modules/vue-loader/lib/runtime/componentNormalizer.js */ "./node_modules/vue-loader/lib/runtime/componentNormalizer.js");





/* normalize component */

var component = Object(_node_modules_vue_loader_lib_runtime_componentNormalizer_js__WEBPACK_IMPORTED_MODULE_2__["default"])(
  _OrderEdit_vue_vue_type_script_lang_js__WEBPACK_IMPORTED_MODULE_1__["default"],
  _OrderEdit_vue_vue_type_template_id_2ac7fed2__WEBPACK_IMPORTED_MODULE_0__["render"],
  _OrderEdit_vue_vue_type_template_id_2ac7fed2__WEBPACK_IMPORTED_MODULE_0__["staticRenderFns"],
  false,
  null,
  null,
  null
  
)

/* hot reload */
if (false) { var api; }
component.options.__file = "resources/js/components/salev2/OrderEdit.vue"
/* harmony default export */ __webpack_exports__["default"] = (component.exports);

/***/ }),

/***/ "./resources/js/components/salev2/OrderEdit.vue?vue&type=script&lang=js":
/*!******************************************************************************!*\
  !*** ./resources/js/components/salev2/OrderEdit.vue?vue&type=script&lang=js ***!
  \******************************************************************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var _node_modules_babel_loader_lib_index_js_ref_4_0_node_modules_vue_loader_lib_index_js_vue_loader_options_OrderEdit_vue_vue_type_script_lang_js__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! -!../../../../node_modules/babel-loader/lib??ref--4-0!../../../../node_modules/vue-loader/lib??vue-loader-options!./OrderEdit.vue?vue&type=script&lang=js */ "./node_modules/babel-loader/lib/index.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/OrderEdit.vue?vue&type=script&lang=js");
/* empty/unused harmony star reexport */ /* harmony default export */ __webpack_exports__["default"] = (_node_modules_babel_loader_lib_index_js_ref_4_0_node_modules_vue_loader_lib_index_js_vue_loader_options_OrderEdit_vue_vue_type_script_lang_js__WEBPACK_IMPORTED_MODULE_0__["default"]); 

/***/ }),

/***/ "./resources/js/components/salev2/OrderEdit.vue?vue&type=template&id=2ac7fed2":
/*!************************************************************************************!*\
  !*** ./resources/js/components/salev2/OrderEdit.vue?vue&type=template&id=2ac7fed2 ***!
  \************************************************************************************/
/*! exports provided: render, staticRenderFns */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var _node_modules_babel_loader_lib_index_js_ref_4_0_node_modules_vue_loader_lib_loaders_templateLoader_js_ref_6_node_modules_vue_loader_lib_index_js_vue_loader_options_OrderEdit_vue_vue_type_template_id_2ac7fed2__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! -!../../../../node_modules/babel-loader/lib??ref--4-0!../../../../node_modules/vue-loader/lib/loaders/templateLoader.js??ref--6!../../../../node_modules/vue-loader/lib??vue-loader-options!./OrderEdit.vue?vue&type=template&id=2ac7fed2 */ "./node_modules/babel-loader/lib/index.js?!./node_modules/vue-loader/lib/loaders/templateLoader.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/OrderEdit.vue?vue&type=template&id=2ac7fed2");
/* harmony reexport (safe) */ __webpack_require__.d(__webpack_exports__, "render", function() { return _node_modules_babel_loader_lib_index_js_ref_4_0_node_modules_vue_loader_lib_loaders_templateLoader_js_ref_6_node_modules_vue_loader_lib_index_js_vue_loader_options_OrderEdit_vue_vue_type_template_id_2ac7fed2__WEBPACK_IMPORTED_MODULE_0__["render"]; });

/* harmony reexport (safe) */ __webpack_require__.d(__webpack_exports__, "staticRenderFns", function() { return _node_modules_babel_loader_lib_index_js_ref_4_0_node_modules_vue_loader_lib_loaders_templateLoader_js_ref_6_node_modules_vue_loader_lib_index_js_vue_loader_options_OrderEdit_vue_vue_type_template_id_2ac7fed2__WEBPACK_IMPORTED_MODULE_0__["staticRenderFns"]; });



/***/ }),

/***/ "./resources/js/components/salev2/OrderManager.vue":
/*!*********************************************************!*\
  !*** ./resources/js/components/salev2/OrderManager.vue ***!
  \*********************************************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var _OrderManager_vue_vue_type_template_id_64110eb0__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ./OrderManager.vue?vue&type=template&id=64110eb0 */ "./resources/js/components/salev2/OrderManager.vue?vue&type=template&id=64110eb0");
/* harmony import */ var _OrderManager_vue_vue_type_script_lang_js__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./OrderManager.vue?vue&type=script&lang=js */ "./resources/js/components/salev2/OrderManager.vue?vue&type=script&lang=js");
/* empty/unused harmony star reexport *//* harmony import */ var _OrderManager_vue_vue_type_style_index_0_id_64110eb0_lang_css__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ./OrderManager.vue?vue&type=style&index=0&id=64110eb0&lang=css */ "./resources/js/components/salev2/OrderManager.vue?vue&type=style&index=0&id=64110eb0&lang=css");
/* harmony import */ var _node_modules_vue_loader_lib_runtime_componentNormalizer_js__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! ../../../../node_modules/vue-loader/lib/runtime/componentNormalizer.js */ "./node_modules/vue-loader/lib/runtime/componentNormalizer.js");






/* normalize component */

var component = Object(_node_modules_vue_loader_lib_runtime_componentNormalizer_js__WEBPACK_IMPORTED_MODULE_3__["default"])(
  _OrderManager_vue_vue_type_script_lang_js__WEBPACK_IMPORTED_MODULE_1__["default"],
  _OrderManager_vue_vue_type_template_id_64110eb0__WEBPACK_IMPORTED_MODULE_0__["render"],
  _OrderManager_vue_vue_type_template_id_64110eb0__WEBPACK_IMPORTED_MODULE_0__["staticRenderFns"],
  false,
  null,
  null,
  null
  
)

/* hot reload */
if (false) { var api; }
component.options.__file = "resources/js/components/salev2/OrderManager.vue"
/* harmony default export */ __webpack_exports__["default"] = (component.exports);

/***/ }),

/***/ "./resources/js/components/salev2/OrderManager.vue?vue&type=script&lang=js":
/*!*********************************************************************************!*\
  !*** ./resources/js/components/salev2/OrderManager.vue?vue&type=script&lang=js ***!
  \*********************************************************************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var _node_modules_babel_loader_lib_index_js_ref_4_0_node_modules_vue_loader_lib_index_js_vue_loader_options_OrderManager_vue_vue_type_script_lang_js__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! -!../../../../node_modules/babel-loader/lib??ref--4-0!../../../../node_modules/vue-loader/lib??vue-loader-options!./OrderManager.vue?vue&type=script&lang=js */ "./node_modules/babel-loader/lib/index.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/OrderManager.vue?vue&type=script&lang=js");
/* empty/unused harmony star reexport */ /* harmony default export */ __webpack_exports__["default"] = (_node_modules_babel_loader_lib_index_js_ref_4_0_node_modules_vue_loader_lib_index_js_vue_loader_options_OrderManager_vue_vue_type_script_lang_js__WEBPACK_IMPORTED_MODULE_0__["default"]); 

/***/ }),

/***/ "./resources/js/components/salev2/OrderManager.vue?vue&type=style&index=0&id=64110eb0&lang=css":
/*!*****************************************************************************************************!*\
  !*** ./resources/js/components/salev2/OrderManager.vue?vue&type=style&index=0&id=64110eb0&lang=css ***!
  \*****************************************************************************************************/
/*! no static exports found */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var _node_modules_style_loader_index_js_node_modules_css_loader_index_js_ref_5_1_node_modules_vue_loader_lib_loaders_stylePostLoader_js_node_modules_postcss_loader_src_index_js_ref_5_2_node_modules_vue_loader_lib_index_js_vue_loader_options_OrderManager_vue_vue_type_style_index_0_id_64110eb0_lang_css__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! -!../../../../node_modules/style-loader!../../../../node_modules/css-loader??ref--5-1!../../../../node_modules/vue-loader/lib/loaders/stylePostLoader.js!../../../../node_modules/postcss-loader/src??ref--5-2!../../../../node_modules/vue-loader/lib??vue-loader-options!./OrderManager.vue?vue&type=style&index=0&id=64110eb0&lang=css */ "./node_modules/style-loader/index.js!./node_modules/css-loader/index.js?!./node_modules/vue-loader/lib/loaders/stylePostLoader.js!./node_modules/postcss-loader/src/index.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/OrderManager.vue?vue&type=style&index=0&id=64110eb0&lang=css");
/* harmony import */ var _node_modules_style_loader_index_js_node_modules_css_loader_index_js_ref_5_1_node_modules_vue_loader_lib_loaders_stylePostLoader_js_node_modules_postcss_loader_src_index_js_ref_5_2_node_modules_vue_loader_lib_index_js_vue_loader_options_OrderManager_vue_vue_type_style_index_0_id_64110eb0_lang_css__WEBPACK_IMPORTED_MODULE_0___default = /*#__PURE__*/__webpack_require__.n(_node_modules_style_loader_index_js_node_modules_css_loader_index_js_ref_5_1_node_modules_vue_loader_lib_loaders_stylePostLoader_js_node_modules_postcss_loader_src_index_js_ref_5_2_node_modules_vue_loader_lib_index_js_vue_loader_options_OrderManager_vue_vue_type_style_index_0_id_64110eb0_lang_css__WEBPACK_IMPORTED_MODULE_0__);
/* harmony reexport (unknown) */ for(var __WEBPACK_IMPORT_KEY__ in _node_modules_style_loader_index_js_node_modules_css_loader_index_js_ref_5_1_node_modules_vue_loader_lib_loaders_stylePostLoader_js_node_modules_postcss_loader_src_index_js_ref_5_2_node_modules_vue_loader_lib_index_js_vue_loader_options_OrderManager_vue_vue_type_style_index_0_id_64110eb0_lang_css__WEBPACK_IMPORTED_MODULE_0__) if(["default"].indexOf(__WEBPACK_IMPORT_KEY__) < 0) (function(key) { __webpack_require__.d(__webpack_exports__, key, function() { return _node_modules_style_loader_index_js_node_modules_css_loader_index_js_ref_5_1_node_modules_vue_loader_lib_loaders_stylePostLoader_js_node_modules_postcss_loader_src_index_js_ref_5_2_node_modules_vue_loader_lib_index_js_vue_loader_options_OrderManager_vue_vue_type_style_index_0_id_64110eb0_lang_css__WEBPACK_IMPORTED_MODULE_0__[key]; }) }(__WEBPACK_IMPORT_KEY__));


/***/ }),

/***/ "./resources/js/components/salev2/OrderManager.vue?vue&type=template&id=64110eb0":
/*!***************************************************************************************!*\
  !*** ./resources/js/components/salev2/OrderManager.vue?vue&type=template&id=64110eb0 ***!
  \***************************************************************************************/
/*! exports provided: render, staticRenderFns */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var _node_modules_babel_loader_lib_index_js_ref_4_0_node_modules_vue_loader_lib_loaders_templateLoader_js_ref_6_node_modules_vue_loader_lib_index_js_vue_loader_options_OrderManager_vue_vue_type_template_id_64110eb0__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! -!../../../../node_modules/babel-loader/lib??ref--4-0!../../../../node_modules/vue-loader/lib/loaders/templateLoader.js??ref--6!../../../../node_modules/vue-loader/lib??vue-loader-options!./OrderManager.vue?vue&type=template&id=64110eb0 */ "./node_modules/babel-loader/lib/index.js?!./node_modules/vue-loader/lib/loaders/templateLoader.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/OrderManager.vue?vue&type=template&id=64110eb0");
/* harmony reexport (safe) */ __webpack_require__.d(__webpack_exports__, "render", function() { return _node_modules_babel_loader_lib_index_js_ref_4_0_node_modules_vue_loader_lib_loaders_templateLoader_js_ref_6_node_modules_vue_loader_lib_index_js_vue_loader_options_OrderManager_vue_vue_type_template_id_64110eb0__WEBPACK_IMPORTED_MODULE_0__["render"]; });

/* harmony reexport (safe) */ __webpack_require__.d(__webpack_exports__, "staticRenderFns", function() { return _node_modules_babel_loader_lib_index_js_ref_4_0_node_modules_vue_loader_lib_loaders_templateLoader_js_ref_6_node_modules_vue_loader_lib_index_js_vue_loader_options_OrderManager_vue_vue_type_template_id_64110eb0__WEBPACK_IMPORTED_MODULE_0__["staticRenderFns"]; });



/***/ }),

/***/ "./resources/js/components/salev2/ProductCreate.vue":
/*!**********************************************************!*\
  !*** ./resources/js/components/salev2/ProductCreate.vue ***!
  \**********************************************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var _ProductCreate_vue_vue_type_template_id_1e3b1bea__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ./ProductCreate.vue?vue&type=template&id=1e3b1bea */ "./resources/js/components/salev2/ProductCreate.vue?vue&type=template&id=1e3b1bea");
/* harmony import */ var _ProductCreate_vue_vue_type_script_lang_js__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./ProductCreate.vue?vue&type=script&lang=js */ "./resources/js/components/salev2/ProductCreate.vue?vue&type=script&lang=js");
/* empty/unused harmony star reexport *//* harmony import */ var _node_modules_vue_loader_lib_runtime_componentNormalizer_js__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ../../../../node_modules/vue-loader/lib/runtime/componentNormalizer.js */ "./node_modules/vue-loader/lib/runtime/componentNormalizer.js");





/* normalize component */

var component = Object(_node_modules_vue_loader_lib_runtime_componentNormalizer_js__WEBPACK_IMPORTED_MODULE_2__["default"])(
  _ProductCreate_vue_vue_type_script_lang_js__WEBPACK_IMPORTED_MODULE_1__["default"],
  _ProductCreate_vue_vue_type_template_id_1e3b1bea__WEBPACK_IMPORTED_MODULE_0__["render"],
  _ProductCreate_vue_vue_type_template_id_1e3b1bea__WEBPACK_IMPORTED_MODULE_0__["staticRenderFns"],
  false,
  null,
  null,
  null
  
)

/* hot reload */
if (false) { var api; }
component.options.__file = "resources/js/components/salev2/ProductCreate.vue"
/* harmony default export */ __webpack_exports__["default"] = (component.exports);

/***/ }),

/***/ "./resources/js/components/salev2/ProductCreate.vue?vue&type=script&lang=js":
/*!**********************************************************************************!*\
  !*** ./resources/js/components/salev2/ProductCreate.vue?vue&type=script&lang=js ***!
  \**********************************************************************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var _node_modules_babel_loader_lib_index_js_ref_4_0_node_modules_vue_loader_lib_index_js_vue_loader_options_ProductCreate_vue_vue_type_script_lang_js__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! -!../../../../node_modules/babel-loader/lib??ref--4-0!../../../../node_modules/vue-loader/lib??vue-loader-options!./ProductCreate.vue?vue&type=script&lang=js */ "./node_modules/babel-loader/lib/index.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/ProductCreate.vue?vue&type=script&lang=js");
/* empty/unused harmony star reexport */ /* harmony default export */ __webpack_exports__["default"] = (_node_modules_babel_loader_lib_index_js_ref_4_0_node_modules_vue_loader_lib_index_js_vue_loader_options_ProductCreate_vue_vue_type_script_lang_js__WEBPACK_IMPORTED_MODULE_0__["default"]); 

/***/ }),

/***/ "./resources/js/components/salev2/ProductCreate.vue?vue&type=template&id=1e3b1bea":
/*!****************************************************************************************!*\
  !*** ./resources/js/components/salev2/ProductCreate.vue?vue&type=template&id=1e3b1bea ***!
  \****************************************************************************************/
/*! exports provided: render, staticRenderFns */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var _node_modules_babel_loader_lib_index_js_ref_4_0_node_modules_vue_loader_lib_loaders_templateLoader_js_ref_6_node_modules_vue_loader_lib_index_js_vue_loader_options_ProductCreate_vue_vue_type_template_id_1e3b1bea__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! -!../../../../node_modules/babel-loader/lib??ref--4-0!../../../../node_modules/vue-loader/lib/loaders/templateLoader.js??ref--6!../../../../node_modules/vue-loader/lib??vue-loader-options!./ProductCreate.vue?vue&type=template&id=1e3b1bea */ "./node_modules/babel-loader/lib/index.js?!./node_modules/vue-loader/lib/loaders/templateLoader.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/ProductCreate.vue?vue&type=template&id=1e3b1bea");
/* harmony reexport (safe) */ __webpack_require__.d(__webpack_exports__, "render", function() { return _node_modules_babel_loader_lib_index_js_ref_4_0_node_modules_vue_loader_lib_loaders_templateLoader_js_ref_6_node_modules_vue_loader_lib_index_js_vue_loader_options_ProductCreate_vue_vue_type_template_id_1e3b1bea__WEBPACK_IMPORTED_MODULE_0__["render"]; });

/* harmony reexport (safe) */ __webpack_require__.d(__webpack_exports__, "staticRenderFns", function() { return _node_modules_babel_loader_lib_index_js_ref_4_0_node_modules_vue_loader_lib_loaders_templateLoader_js_ref_6_node_modules_vue_loader_lib_index_js_vue_loader_options_ProductCreate_vue_vue_type_template_id_1e3b1bea__WEBPACK_IMPORTED_MODULE_0__["staticRenderFns"]; });



/***/ }),

/***/ "./resources/js/components/salev2/ProductEdit.vue":
/*!********************************************************!*\
  !*** ./resources/js/components/salev2/ProductEdit.vue ***!
  \********************************************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var _ProductEdit_vue_vue_type_template_id_635d4f78__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ./ProductEdit.vue?vue&type=template&id=635d4f78 */ "./resources/js/components/salev2/ProductEdit.vue?vue&type=template&id=635d4f78");
/* harmony import */ var _ProductEdit_vue_vue_type_script_lang_js__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./ProductEdit.vue?vue&type=script&lang=js */ "./resources/js/components/salev2/ProductEdit.vue?vue&type=script&lang=js");
/* empty/unused harmony star reexport *//* harmony import */ var _node_modules_vue_loader_lib_runtime_componentNormalizer_js__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ../../../../node_modules/vue-loader/lib/runtime/componentNormalizer.js */ "./node_modules/vue-loader/lib/runtime/componentNormalizer.js");





/* normalize component */

var component = Object(_node_modules_vue_loader_lib_runtime_componentNormalizer_js__WEBPACK_IMPORTED_MODULE_2__["default"])(
  _ProductEdit_vue_vue_type_script_lang_js__WEBPACK_IMPORTED_MODULE_1__["default"],
  _ProductEdit_vue_vue_type_template_id_635d4f78__WEBPACK_IMPORTED_MODULE_0__["render"],
  _ProductEdit_vue_vue_type_template_id_635d4f78__WEBPACK_IMPORTED_MODULE_0__["staticRenderFns"],
  false,
  null,
  null,
  null
  
)

/* hot reload */
if (false) { var api; }
component.options.__file = "resources/js/components/salev2/ProductEdit.vue"
/* harmony default export */ __webpack_exports__["default"] = (component.exports);

/***/ }),

/***/ "./resources/js/components/salev2/ProductEdit.vue?vue&type=script&lang=js":
/*!********************************************************************************!*\
  !*** ./resources/js/components/salev2/ProductEdit.vue?vue&type=script&lang=js ***!
  \********************************************************************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var _node_modules_babel_loader_lib_index_js_ref_4_0_node_modules_vue_loader_lib_index_js_vue_loader_options_ProductEdit_vue_vue_type_script_lang_js__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! -!../../../../node_modules/babel-loader/lib??ref--4-0!../../../../node_modules/vue-loader/lib??vue-loader-options!./ProductEdit.vue?vue&type=script&lang=js */ "./node_modules/babel-loader/lib/index.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/ProductEdit.vue?vue&type=script&lang=js");
/* empty/unused harmony star reexport */ /* harmony default export */ __webpack_exports__["default"] = (_node_modules_babel_loader_lib_index_js_ref_4_0_node_modules_vue_loader_lib_index_js_vue_loader_options_ProductEdit_vue_vue_type_script_lang_js__WEBPACK_IMPORTED_MODULE_0__["default"]); 

/***/ }),

/***/ "./resources/js/components/salev2/ProductEdit.vue?vue&type=template&id=635d4f78":
/*!**************************************************************************************!*\
  !*** ./resources/js/components/salev2/ProductEdit.vue?vue&type=template&id=635d4f78 ***!
  \**************************************************************************************/
/*! exports provided: render, staticRenderFns */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var _node_modules_babel_loader_lib_index_js_ref_4_0_node_modules_vue_loader_lib_loaders_templateLoader_js_ref_6_node_modules_vue_loader_lib_index_js_vue_loader_options_ProductEdit_vue_vue_type_template_id_635d4f78__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! -!../../../../node_modules/babel-loader/lib??ref--4-0!../../../../node_modules/vue-loader/lib/loaders/templateLoader.js??ref--6!../../../../node_modules/vue-loader/lib??vue-loader-options!./ProductEdit.vue?vue&type=template&id=635d4f78 */ "./node_modules/babel-loader/lib/index.js?!./node_modules/vue-loader/lib/loaders/templateLoader.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/ProductEdit.vue?vue&type=template&id=635d4f78");
/* harmony reexport (safe) */ __webpack_require__.d(__webpack_exports__, "render", function() { return _node_modules_babel_loader_lib_index_js_ref_4_0_node_modules_vue_loader_lib_loaders_templateLoader_js_ref_6_node_modules_vue_loader_lib_index_js_vue_loader_options_ProductEdit_vue_vue_type_template_id_635d4f78__WEBPACK_IMPORTED_MODULE_0__["render"]; });

/* harmony reexport (safe) */ __webpack_require__.d(__webpack_exports__, "staticRenderFns", function() { return _node_modules_babel_loader_lib_index_js_ref_4_0_node_modules_vue_loader_lib_loaders_templateLoader_js_ref_6_node_modules_vue_loader_lib_index_js_vue_loader_options_ProductEdit_vue_vue_type_template_id_635d4f78__WEBPACK_IMPORTED_MODULE_0__["staticRenderFns"]; });



/***/ }),

/***/ "./resources/js/components/salev2/ProductVault.vue":
/*!*********************************************************!*\
  !*** ./resources/js/components/salev2/ProductVault.vue ***!
  \*********************************************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var _ProductVault_vue_vue_type_template_id_1a1232b4__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ./ProductVault.vue?vue&type=template&id=1a1232b4 */ "./resources/js/components/salev2/ProductVault.vue?vue&type=template&id=1a1232b4");
/* harmony import */ var _ProductVault_vue_vue_type_script_lang_js__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./ProductVault.vue?vue&type=script&lang=js */ "./resources/js/components/salev2/ProductVault.vue?vue&type=script&lang=js");
/* empty/unused harmony star reexport *//* harmony import */ var _ProductVault_vue_vue_type_style_index_0_id_1a1232b4_lang_css__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ./ProductVault.vue?vue&type=style&index=0&id=1a1232b4&lang=css */ "./resources/js/components/salev2/ProductVault.vue?vue&type=style&index=0&id=1a1232b4&lang=css");
/* harmony import */ var _node_modules_vue_loader_lib_runtime_componentNormalizer_js__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! ../../../../node_modules/vue-loader/lib/runtime/componentNormalizer.js */ "./node_modules/vue-loader/lib/runtime/componentNormalizer.js");






/* normalize component */

var component = Object(_node_modules_vue_loader_lib_runtime_componentNormalizer_js__WEBPACK_IMPORTED_MODULE_3__["default"])(
  _ProductVault_vue_vue_type_script_lang_js__WEBPACK_IMPORTED_MODULE_1__["default"],
  _ProductVault_vue_vue_type_template_id_1a1232b4__WEBPACK_IMPORTED_MODULE_0__["render"],
  _ProductVault_vue_vue_type_template_id_1a1232b4__WEBPACK_IMPORTED_MODULE_0__["staticRenderFns"],
  false,
  null,
  null,
  null
  
)

/* hot reload */
if (false) { var api; }
component.options.__file = "resources/js/components/salev2/ProductVault.vue"
/* harmony default export */ __webpack_exports__["default"] = (component.exports);

/***/ }),

/***/ "./resources/js/components/salev2/ProductVault.vue?vue&type=script&lang=js":
/*!*********************************************************************************!*\
  !*** ./resources/js/components/salev2/ProductVault.vue?vue&type=script&lang=js ***!
  \*********************************************************************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var _node_modules_babel_loader_lib_index_js_ref_4_0_node_modules_vue_loader_lib_index_js_vue_loader_options_ProductVault_vue_vue_type_script_lang_js__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! -!../../../../node_modules/babel-loader/lib??ref--4-0!../../../../node_modules/vue-loader/lib??vue-loader-options!./ProductVault.vue?vue&type=script&lang=js */ "./node_modules/babel-loader/lib/index.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/ProductVault.vue?vue&type=script&lang=js");
/* empty/unused harmony star reexport */ /* harmony default export */ __webpack_exports__["default"] = (_node_modules_babel_loader_lib_index_js_ref_4_0_node_modules_vue_loader_lib_index_js_vue_loader_options_ProductVault_vue_vue_type_script_lang_js__WEBPACK_IMPORTED_MODULE_0__["default"]); 

/***/ }),

/***/ "./resources/js/components/salev2/ProductVault.vue?vue&type=style&index=0&id=1a1232b4&lang=css":
/*!*****************************************************************************************************!*\
  !*** ./resources/js/components/salev2/ProductVault.vue?vue&type=style&index=0&id=1a1232b4&lang=css ***!
  \*****************************************************************************************************/
/*! no static exports found */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var _node_modules_style_loader_index_js_node_modules_css_loader_index_js_ref_5_1_node_modules_vue_loader_lib_loaders_stylePostLoader_js_node_modules_postcss_loader_src_index_js_ref_5_2_node_modules_vue_loader_lib_index_js_vue_loader_options_ProductVault_vue_vue_type_style_index_0_id_1a1232b4_lang_css__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! -!../../../../node_modules/style-loader!../../../../node_modules/css-loader??ref--5-1!../../../../node_modules/vue-loader/lib/loaders/stylePostLoader.js!../../../../node_modules/postcss-loader/src??ref--5-2!../../../../node_modules/vue-loader/lib??vue-loader-options!./ProductVault.vue?vue&type=style&index=0&id=1a1232b4&lang=css */ "./node_modules/style-loader/index.js!./node_modules/css-loader/index.js?!./node_modules/vue-loader/lib/loaders/stylePostLoader.js!./node_modules/postcss-loader/src/index.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/ProductVault.vue?vue&type=style&index=0&id=1a1232b4&lang=css");
/* harmony import */ var _node_modules_style_loader_index_js_node_modules_css_loader_index_js_ref_5_1_node_modules_vue_loader_lib_loaders_stylePostLoader_js_node_modules_postcss_loader_src_index_js_ref_5_2_node_modules_vue_loader_lib_index_js_vue_loader_options_ProductVault_vue_vue_type_style_index_0_id_1a1232b4_lang_css__WEBPACK_IMPORTED_MODULE_0___default = /*#__PURE__*/__webpack_require__.n(_node_modules_style_loader_index_js_node_modules_css_loader_index_js_ref_5_1_node_modules_vue_loader_lib_loaders_stylePostLoader_js_node_modules_postcss_loader_src_index_js_ref_5_2_node_modules_vue_loader_lib_index_js_vue_loader_options_ProductVault_vue_vue_type_style_index_0_id_1a1232b4_lang_css__WEBPACK_IMPORTED_MODULE_0__);
/* harmony reexport (unknown) */ for(var __WEBPACK_IMPORT_KEY__ in _node_modules_style_loader_index_js_node_modules_css_loader_index_js_ref_5_1_node_modules_vue_loader_lib_loaders_stylePostLoader_js_node_modules_postcss_loader_src_index_js_ref_5_2_node_modules_vue_loader_lib_index_js_vue_loader_options_ProductVault_vue_vue_type_style_index_0_id_1a1232b4_lang_css__WEBPACK_IMPORTED_MODULE_0__) if(["default"].indexOf(__WEBPACK_IMPORT_KEY__) < 0) (function(key) { __webpack_require__.d(__webpack_exports__, key, function() { return _node_modules_style_loader_index_js_node_modules_css_loader_index_js_ref_5_1_node_modules_vue_loader_lib_loaders_stylePostLoader_js_node_modules_postcss_loader_src_index_js_ref_5_2_node_modules_vue_loader_lib_index_js_vue_loader_options_ProductVault_vue_vue_type_style_index_0_id_1a1232b4_lang_css__WEBPACK_IMPORTED_MODULE_0__[key]; }) }(__WEBPACK_IMPORT_KEY__));


/***/ }),

/***/ "./resources/js/components/salev2/ProductVault.vue?vue&type=template&id=1a1232b4":
/*!***************************************************************************************!*\
  !*** ./resources/js/components/salev2/ProductVault.vue?vue&type=template&id=1a1232b4 ***!
  \***************************************************************************************/
/*! exports provided: render, staticRenderFns */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var _node_modules_babel_loader_lib_index_js_ref_4_0_node_modules_vue_loader_lib_loaders_templateLoader_js_ref_6_node_modules_vue_loader_lib_index_js_vue_loader_options_ProductVault_vue_vue_type_template_id_1a1232b4__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! -!../../../../node_modules/babel-loader/lib??ref--4-0!../../../../node_modules/vue-loader/lib/loaders/templateLoader.js??ref--6!../../../../node_modules/vue-loader/lib??vue-loader-options!./ProductVault.vue?vue&type=template&id=1a1232b4 */ "./node_modules/babel-loader/lib/index.js?!./node_modules/vue-loader/lib/loaders/templateLoader.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/ProductVault.vue?vue&type=template&id=1a1232b4");
/* harmony reexport (safe) */ __webpack_require__.d(__webpack_exports__, "render", function() { return _node_modules_babel_loader_lib_index_js_ref_4_0_node_modules_vue_loader_lib_loaders_templateLoader_js_ref_6_node_modules_vue_loader_lib_index_js_vue_loader_options_ProductVault_vue_vue_type_template_id_1a1232b4__WEBPACK_IMPORTED_MODULE_0__["render"]; });

/* harmony reexport (safe) */ __webpack_require__.d(__webpack_exports__, "staticRenderFns", function() { return _node_modules_babel_loader_lib_index_js_ref_4_0_node_modules_vue_loader_lib_loaders_templateLoader_js_ref_6_node_modules_vue_loader_lib_index_js_vue_loader_options_ProductVault_vue_vue_type_template_id_1a1232b4__WEBPACK_IMPORTED_MODULE_0__["staticRenderFns"]; });



/***/ }),

/***/ "./resources/js/components/salev2/SmartUpload.vue":
/*!********************************************************!*\
  !*** ./resources/js/components/salev2/SmartUpload.vue ***!
  \********************************************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var _SmartUpload_vue_vue_type_template_id_30d6d489__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ./SmartUpload.vue?vue&type=template&id=30d6d489 */ "./resources/js/components/salev2/SmartUpload.vue?vue&type=template&id=30d6d489");
/* harmony import */ var _SmartUpload_vue_vue_type_script_lang_js__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./SmartUpload.vue?vue&type=script&lang=js */ "./resources/js/components/salev2/SmartUpload.vue?vue&type=script&lang=js");
/* empty/unused harmony star reexport *//* harmony import */ var _SmartUpload_vue_vue_type_style_index_0_id_30d6d489_lang_css__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ./SmartUpload.vue?vue&type=style&index=0&id=30d6d489&lang=css */ "./resources/js/components/salev2/SmartUpload.vue?vue&type=style&index=0&id=30d6d489&lang=css");
/* harmony import */ var _node_modules_vue_loader_lib_runtime_componentNormalizer_js__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! ../../../../node_modules/vue-loader/lib/runtime/componentNormalizer.js */ "./node_modules/vue-loader/lib/runtime/componentNormalizer.js");






/* normalize component */

var component = Object(_node_modules_vue_loader_lib_runtime_componentNormalizer_js__WEBPACK_IMPORTED_MODULE_3__["default"])(
  _SmartUpload_vue_vue_type_script_lang_js__WEBPACK_IMPORTED_MODULE_1__["default"],
  _SmartUpload_vue_vue_type_template_id_30d6d489__WEBPACK_IMPORTED_MODULE_0__["render"],
  _SmartUpload_vue_vue_type_template_id_30d6d489__WEBPACK_IMPORTED_MODULE_0__["staticRenderFns"],
  false,
  null,
  null,
  null
  
)

/* hot reload */
if (false) { var api; }
component.options.__file = "resources/js/components/salev2/SmartUpload.vue"
/* harmony default export */ __webpack_exports__["default"] = (component.exports);

/***/ }),

/***/ "./resources/js/components/salev2/SmartUpload.vue?vue&type=script&lang=js":
/*!********************************************************************************!*\
  !*** ./resources/js/components/salev2/SmartUpload.vue?vue&type=script&lang=js ***!
  \********************************************************************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var _node_modules_babel_loader_lib_index_js_ref_4_0_node_modules_vue_loader_lib_index_js_vue_loader_options_SmartUpload_vue_vue_type_script_lang_js__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! -!../../../../node_modules/babel-loader/lib??ref--4-0!../../../../node_modules/vue-loader/lib??vue-loader-options!./SmartUpload.vue?vue&type=script&lang=js */ "./node_modules/babel-loader/lib/index.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/SmartUpload.vue?vue&type=script&lang=js");
/* empty/unused harmony star reexport */ /* harmony default export */ __webpack_exports__["default"] = (_node_modules_babel_loader_lib_index_js_ref_4_0_node_modules_vue_loader_lib_index_js_vue_loader_options_SmartUpload_vue_vue_type_script_lang_js__WEBPACK_IMPORTED_MODULE_0__["default"]); 

/***/ }),

/***/ "./resources/js/components/salev2/SmartUpload.vue?vue&type=style&index=0&id=30d6d489&lang=css":
/*!****************************************************************************************************!*\
  !*** ./resources/js/components/salev2/SmartUpload.vue?vue&type=style&index=0&id=30d6d489&lang=css ***!
  \****************************************************************************************************/
/*! no static exports found */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var _node_modules_style_loader_index_js_node_modules_css_loader_index_js_ref_5_1_node_modules_vue_loader_lib_loaders_stylePostLoader_js_node_modules_postcss_loader_src_index_js_ref_5_2_node_modules_vue_loader_lib_index_js_vue_loader_options_SmartUpload_vue_vue_type_style_index_0_id_30d6d489_lang_css__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! -!../../../../node_modules/style-loader!../../../../node_modules/css-loader??ref--5-1!../../../../node_modules/vue-loader/lib/loaders/stylePostLoader.js!../../../../node_modules/postcss-loader/src??ref--5-2!../../../../node_modules/vue-loader/lib??vue-loader-options!./SmartUpload.vue?vue&type=style&index=0&id=30d6d489&lang=css */ "./node_modules/style-loader/index.js!./node_modules/css-loader/index.js?!./node_modules/vue-loader/lib/loaders/stylePostLoader.js!./node_modules/postcss-loader/src/index.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/SmartUpload.vue?vue&type=style&index=0&id=30d6d489&lang=css");
/* harmony import */ var _node_modules_style_loader_index_js_node_modules_css_loader_index_js_ref_5_1_node_modules_vue_loader_lib_loaders_stylePostLoader_js_node_modules_postcss_loader_src_index_js_ref_5_2_node_modules_vue_loader_lib_index_js_vue_loader_options_SmartUpload_vue_vue_type_style_index_0_id_30d6d489_lang_css__WEBPACK_IMPORTED_MODULE_0___default = /*#__PURE__*/__webpack_require__.n(_node_modules_style_loader_index_js_node_modules_css_loader_index_js_ref_5_1_node_modules_vue_loader_lib_loaders_stylePostLoader_js_node_modules_postcss_loader_src_index_js_ref_5_2_node_modules_vue_loader_lib_index_js_vue_loader_options_SmartUpload_vue_vue_type_style_index_0_id_30d6d489_lang_css__WEBPACK_IMPORTED_MODULE_0__);
/* harmony reexport (unknown) */ for(var __WEBPACK_IMPORT_KEY__ in _node_modules_style_loader_index_js_node_modules_css_loader_index_js_ref_5_1_node_modules_vue_loader_lib_loaders_stylePostLoader_js_node_modules_postcss_loader_src_index_js_ref_5_2_node_modules_vue_loader_lib_index_js_vue_loader_options_SmartUpload_vue_vue_type_style_index_0_id_30d6d489_lang_css__WEBPACK_IMPORTED_MODULE_0__) if(["default"].indexOf(__WEBPACK_IMPORT_KEY__) < 0) (function(key) { __webpack_require__.d(__webpack_exports__, key, function() { return _node_modules_style_loader_index_js_node_modules_css_loader_index_js_ref_5_1_node_modules_vue_loader_lib_loaders_stylePostLoader_js_node_modules_postcss_loader_src_index_js_ref_5_2_node_modules_vue_loader_lib_index_js_vue_loader_options_SmartUpload_vue_vue_type_style_index_0_id_30d6d489_lang_css__WEBPACK_IMPORTED_MODULE_0__[key]; }) }(__WEBPACK_IMPORT_KEY__));


/***/ }),

/***/ "./resources/js/components/salev2/SmartUpload.vue?vue&type=template&id=30d6d489":
/*!**************************************************************************************!*\
  !*** ./resources/js/components/salev2/SmartUpload.vue?vue&type=template&id=30d6d489 ***!
  \**************************************************************************************/
/*! exports provided: render, staticRenderFns */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var _node_modules_babel_loader_lib_index_js_ref_4_0_node_modules_vue_loader_lib_loaders_templateLoader_js_ref_6_node_modules_vue_loader_lib_index_js_vue_loader_options_SmartUpload_vue_vue_type_template_id_30d6d489__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! -!../../../../node_modules/babel-loader/lib??ref--4-0!../../../../node_modules/vue-loader/lib/loaders/templateLoader.js??ref--6!../../../../node_modules/vue-loader/lib??vue-loader-options!./SmartUpload.vue?vue&type=template&id=30d6d489 */ "./node_modules/babel-loader/lib/index.js?!./node_modules/vue-loader/lib/loaders/templateLoader.js?!./node_modules/vue-loader/lib/index.js?!./resources/js/components/salev2/SmartUpload.vue?vue&type=template&id=30d6d489");
/* harmony reexport (safe) */ __webpack_require__.d(__webpack_exports__, "render", function() { return _node_modules_babel_loader_lib_index_js_ref_4_0_node_modules_vue_loader_lib_loaders_templateLoader_js_ref_6_node_modules_vue_loader_lib_index_js_vue_loader_options_SmartUpload_vue_vue_type_template_id_30d6d489__WEBPACK_IMPORTED_MODULE_0__["render"]; });

/* harmony reexport (safe) */ __webpack_require__.d(__webpack_exports__, "staticRenderFns", function() { return _node_modules_babel_loader_lib_index_js_ref_4_0_node_modules_vue_loader_lib_loaders_templateLoader_js_ref_6_node_modules_vue_loader_lib_index_js_vue_loader_options_SmartUpload_vue_vue_type_template_id_30d6d489__WEBPACK_IMPORTED_MODULE_0__["staticRenderFns"]; });



/***/ }),

/***/ 0:
/*!***********************************!*\
  !*** multi ./resources/js/app.js ***!
  \***********************************/
/*! no static exports found */
/***/ (function(module, exports, __webpack_require__) {

module.exports = __webpack_require__(/*! /Users/macbook/Data/Source/vintagealwaystrue/resources/js/app.js */"./resources/js/app.js");


/***/ })

/******/ });