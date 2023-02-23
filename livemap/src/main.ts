import "./style.css";
import { Hash } from "./hash";
import L, { LatLng } from "leaflet";

const center_x = 117.3;
const center_y = 172.8;
const scale_x = 0.02072;
const scale_y = 0.0205;

let CUSTOM_CRS = L.extend({}, L.CRS.Simple, {
  projection: L.Projection.LonLat,
  scale: function (zoom: number) {
    return Math.pow(2, zoom);
  },
  zoom: function (sc: number) {
    return Math.log(sc) / 0.6931471805599453;
  },
  distance: function (pos1: LatLng, pos2: LatLng) {
    var x_difference = pos2.lng - pos1.lng;
    var y_difference = pos2.lat - pos1.lat;
    return Math.sqrt(x_difference * x_difference + y_difference * y_difference);
  },
  transformation: new L.Transformation(scale_x, center_x, -scale_y, center_y),
  infinite: true,
});

// Create layers
var atlas = L.tileLayer("tiles/atlas/{z}/{x}/{y}.png", {
  attribution:
    '<a href="http://www.rockstargames.com/V/">Grand Theft Auto V</a>, web version quickly done by <a href="http://www.somebits.com/weblog/">Nelson Minar</a>',
  minZoom: 1,
  maxZoom: 6,
  noWrap: false,
  tms: true,
});
var road = L.tileLayer("tiles/road/{z}/{x}/{y}.png", {
  attribution:
    '<a href="http://www.rockstargames.com/V/">Grand Theft Auto V</a>, web version quickly done by <a href="http://www.somebits.com/weblog/">Nelson Minar</a>',
  minZoom: 1,
  maxZoom: 6,
  noWrap: false,
  tms: true,
});
var satelite = L.tileLayer("tiles/satelite/{z}/{x}/{y}.png", {
  attribution:
    '<a href="http://www.rockstargames.com/V/">Grand Theft Auto V</a>, web version quickly done by <a href="http://www.somebits.com/weblog/">Nelson Minar</a>',
  minZoom: 1,
  maxZoom: 6,
  noWrap: false,
  tms: true,
});

var map = L.map("map", { layers: [satelite], crs: CUSTOM_CRS });
new Hash(map);
map.setView([0, 0], 2);
let markersLayer = new L.LayerGroup().addTo(map);
L.control
  .layers({ Satelite: satelite, Atlas: atlas, Road: road }, { Markers: markersLayer })
  .addTo(map);
satelite.bringToFront();

// Latitude and Longitiude popup on mouse over
var _latlng: HTMLDivElement;
let Position = L.Control.extend({
  _container: null,
  options: {
    position: "bottomleft",
  },
  onAdd: function () {
    var latlng = L.DomUtil.create("div", "mouseposition");
    _latlng = latlng;
    return latlng;
  },
  updateHTML: function (lat: number, lng: number) {
    _latlng.innerHTML = "Latitude: " + lat + "   Longitiude: " + lng;
  },
});

var position = new Position();
map.addControl(position);

map.addEventListener("mousemove", (event) => {
  let lat = Math.round(event.latlng.lat * 100000) / 100000;
  let lng = Math.round(event.latlng.lng * 100000) / 100000;
  position.updateHTML(lat, lng);
});
