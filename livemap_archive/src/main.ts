import './style.css';

import L from 'leaflet';

import { Livemap } from './class/Livemap';
import { customCRS } from './config/CRS';

const atlas = L.tileLayer('tiles/atlas/{z}/{x}/{y}.png', {
	attribution:
		'<a href="http://www.rockstargames.com/V/">Grand Theft Auto V</a>, web version quickly done by <a href="http://www.somebits.com/weblog/">Nelson Minar</a>',
	minZoom: 1,
	maxZoom: 6,
	noWrap: false,
	tms: true,
});
const road = L.tileLayer('tiles/road/{z}/{x}/{y}.png', {
	attribution:
		'<a href="http://www.rockstargames.com/V/">Grand Theft Auto V</a>, web version quickly done by <a href="http://www.somebits.com/weblog/">Nelson Minar</a>',
	minZoom: 1,
	maxZoom: 6,
	noWrap: false,
	tms: true,
});
const satelite = L.tileLayer('tiles/satelite/{z}/{x}/{y}.png', {
	attribution:
		'<a href="http://www.rockstargames.com/V/">Grand Theft Auto V</a>, web version quickly done by <a href="http://www.somebits.com/weblog/">Nelson Minar</a>',
	minZoom: 1,
	maxZoom: 6,
	noWrap: false,
	tms: true,
});

const map = new Livemap('map', { layers: [satelite], crs: customCRS });
map.addHash();
map.setView([0, 0], 2);

const markersLayer = new L.LayerGroup().addTo(map);
L.control.layers({ Satelite: satelite, Atlas: atlas, Road: road }, { Markers: markersLayer }).addTo(map);
satelite.bringToFront();

// Latitude and Longitiude popup on mouse over
let _latlng: HTMLDivElement;
const Position = L.Control.extend({
	_container: null,
	options: {
		position: 'bottomleft',
	},
	onAdd: function () {
		const latlng = L.DomUtil.create('div', 'mouseposition');
		_latlng = latlng;
		return latlng;
	},
	updateHTML: function (lat: number, lng: number) {
		_latlng.innerHTML = 'Latitude: ' + lat + '   Longitiude: ' + lng;
	},
});

const position = new Position();
map.addControl(position);

map.addEventListener('mousemove', (event) => {
	const lat = Math.round(event.latlng.lat * 100000) / 100000;
	const lng = Math.round(event.latlng.lng * 100000) / 100000;
	position.updateHTML(lat, lng);
});
