<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Chronometer</title>
</head>
<body>

<label for="chrono">Chronometer:</label><input type="number" id="chrono" value="{{.Chrono}}" disabled><br/>
<label for="point">Checkpoint:</label><input type="number" id="point" value="{{.Point}}" disabled><br/>
<label for="number">Number:</label><input type="number" id="number"><br/>

<button type="button" onclick="chrono_AppendPass();">Append</button>
<button type="button" onclick="if (confirm('Undo last pass?')) chrono_UndoPass();">Undo</button>
<button type="button" onclick="if (confirm('Sync passes?')) chrono_SyncPasses();">Sync</button>
<a href="/copy?c={{.Chrono}}&p={{.Point}}" download>Download</a><br />

<div id="console"></div>

<script>

const g_Chrono = {{.Chrono}};
const g_Point = {{.Point}};

const g_Address = '/sync?c=' + g_Chrono + '&p=' + g_Point;
const g_KeyPrefix = 'chrono-' + g_Chrono + ':point-' + g_Point;
const g_FirstKey = g_KeyPrefix + ':first';
const g_LastKey = g_KeyPrefix + ':last';

var g_Number = document.getElementById('number');
var g_Console = document.getElementById('console');
var g_Storage = window.localStorage;

function chrono_Now() {
	return new Date().getTime();
}

function chrono_Log(s) {
	var d = new Date();
	var t = d.toLocaleTimeString();
	var e = document.createElement('div');
	e.innerHTML = '[' + t + '] ' + s;
	if (g_Console.childNodes.length > 0) {
		g_Console.insertBefore(e, g_Console.firstChild);
	} else {
		g_Console.appendChild(e);
	}
}

function chrono_PassKey(i) {
	return g_KeyPrefix + ':pass-' + i;
}

function chrono_Incr(k) {
	var i = g_Storage.getItem(k);
	if (i == null) {
		i = 1000;
	}
	i++;
	g_Storage.setItem(k, i);
	return i;
}

function chrono_AppendPass() {
	var n = g_Number.value;
	if (n == '') {
		chrono_Log('<b>Error</b>: Number not entered!');
		return;
	}
        var t = chrono_Now();
        var i = chrono_Incr(g_LastKey);
        var k = chrono_PassKey(i);
	g_Storage.setItem(k, n + ',' + t);
	chrono_Log('Pass ' + i + ' number <b>' + n + '</b> time ' + t);
}

function chrono_UndoPass() {
        var a = g_Storage.getItem(g_FirstKey);
        if (a == null) {
                a = 0;
        }
        var b = g_Storage.getItem(g_LastKey);
        if (b == null) {
                chrono_Log('<b>Error</b>: Last index unknown!');
                return;
        }
	if (a == b) {
		chrono_Log('<b>Warning</b>: All passes synced!');
                return;
	}
	var k = chrono_PassKey(b);
	var v = g_Storage.getItem(k);
	if (v == null) {
		chrono_Log('<b>Error</b>: Pass not exists!');
                return;
	}
        if (v == '') {
                chrono_Log('<b>Warning</b>: Pass already undone!');
                return;
        }
	g_Storage.setItem(k, '');
	chrono_Log('Pass ' + b + ' undone');
}

function chrono_PostDataCallback(x, a, b) {
	if (x.readyState == 4 && x.status == 200) {
		g_Storage.setItem(g_FirstKey, b)
		chrono_Log('Sync from ' + a + ' to ' + b)
	}
}

function chrono_SyncPasses() {
	var a = g_Storage.getItem(g_FirstKey);
	if (a == null) {
		a = 0;
	}
	var b = g_Storage.getItem(g_LastKey);
	if (b == null) {
		chrono_Log('<b>Error</b>: Last index unknown!');
		return;
	}
	var p = b;
	var d = '';
	var k = chrono_PassKey(b);
	var v = g_Storage.getItem(k);
	while ((v != null) && (b > a)) {
		d = d + v + '\n';
		b--;
		k = chrono_PassKey(b);
		v = g_Storage.getItem(k);
	}
	if (d == '') {
		chrono_Log('<b>Warning</b>: All passes synced!');
		return;
	}
	var x = new XMLHttpRequest();
	x.timeout = 2000;
	x.onreadystatechange = function() {
		chrono_PostDataCallback(x, b + 1, p);
	};
	x.open('POST', g_Address, true);
	x.send(d);
}

</script>

</body>
</html>

