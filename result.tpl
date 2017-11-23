<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Chronometer</title>
<style>

tr:nth-child(even) { background: #eee; }
tr:nth-child(odd) { background: #fff; }

table { counter-reset: rowNumber; }
table tr > td:first-child::before { counter-increment: rowNumber; }
table tr td:first-child::before {
	content: counter(rowNumber);
	min-width: 1em;
	margin-right: 0.5em;
}

</style>
</head>
<body>
<table>
<tr>
<th>#</th>
<th>Number</th>
<th>Name</th>
<th>Team</th>
<th>Category</th>
<th>Laps</th>
<th>Duration</th>
<th>Lap 1</th>
<th>Lap 2</th>
<th>Lap 3</th>
<th>Lap 4</th>
<th>Lap 5</th>
</tr>
{{range .}}
<tr>
<td></td>
<td>{{.Number}}</td>
<td>{{.Name}}</td>
<td>{{.Team}}</td>
<td>{{.Category}}</td>
<td>{{.Laps}}</td>
<td>{{.Dur}}</td>
<td>{{.Dur1}}</td>
<td>{{.Dur2}}</td>
<td>{{.Dur3}}</td>
<td>{{.Dur4}}</td>
<td>{{.Dur5}}</td>
</tr>
{{end}}
</table>
</body>
</html>

