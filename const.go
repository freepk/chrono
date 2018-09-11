package main

const (
	sqlCreateDb = `CREATE TABLE IF NOT EXISTS Participants (number INTEGER NOT NULL, name TEXT NOT NULL, category TEXT NOT NULL, team TEXT, CONSTRAINT PK_Participants PRIMARY KEY(number));
WITH Nums AS (
	SELECT 0 n
	UNION ALL SELECT 1
	UNION ALL SELECT 2
	UNION ALL SELECT 3
	UNION ALL SELECT 4
	UNION ALL SELECT 5
	UNION ALL SELECT 6
	UNION ALL SELECT 7
	UNION ALL SELECT 8
	UNION ALL SELECT 9
)
INSERT INTO Participants(number, name, category, team)
SELECT n0.n + n1.n * 10 + 1 as number, '' name, '' category, '' team
FROM Nums n0
	CROSS JOIN Nums n1;
CREATE TABLE IF NOT EXISTS Passes (number INT NOT NULL, point INT NOT NULL, pass INT NOT NULL);
CREATE INDEX IF NOT EXISTS IX_Passes_number_pass ON Passes (number, pass);`
	sqlResults = `WITH FirstPasses AS (
	SELECT DISTINCT number, pass
	FROM Passes p0
	WHERE NOT EXISTS (SELECT 1 FROM Passes p1 WHERE (p1.number = p0.number) AND (p1.pass < p0.pass))
), ValidPasses AS (
	SELECT DISTINCT number, pass
	FROM Passes p0
	WHERE NOT EXISTS (SELECT 1 FROM Passes p1 WHERE (p1.number = p0.number) AND (p1.pass < p0.pass) AND (p0.pass - p1.pass < 300000))
)
, Level1Passes AS (SELECT p0.*, (SELECT MIN(p1.pass) FROM ValidPasses p1 WHERE (p1.number = p0.number) AND (p1.pass > p0.pass )) pass1 FROM FirstPasses  p0)
, Level2Passes AS (SELECT p0.*, (SELECT MIN(p1.pass) FROM ValidPasses p1 WHERE (p1.number = p0.number) AND (p1.pass > p0.pass1)) pass2 FROM Level1Passes p0)
, Level3Passes AS (SELECT p0.*, (SELECT MIN(p1.pass) FROM ValidPasses p1 WHERE (p1.number = p0.number) AND (p1.pass > p0.pass2)) pass3 FROM Level2Passes p0)
, Level4Passes AS (SELECT p0.*, (SELECT MIN(p1.pass) FROM ValidPasses p1 WHERE (p1.number = p0.number) AND (p1.pass > p0.pass3)) pass4 FROM Level3Passes p0)
, Level5Passes AS (SELECT p0.*, (SELECT MIN(p1.pass) FROM ValidPasses p1 WHERE (p1.number = p0.number) AND (p1.pass > p0.pass4)) pass5 FROM Level4Passes p0)
, Results AS (
	SELECT number
		, CASE
			WHEN pass5 IS NOT NULL THEN 5
			WHEN pass4 IS NOT NULL THEN 4
			WHEN pass3 IS NOT NULL THEN 3
			WHEN pass2 IS NOT NULL THEN 2
			WHEN pass1 IS NOT NULL THEN 1
		END laps
		, CASE
			--WHEN pass5 IS NOT NULL THEN pass5
			--WHEN pass4 IS NOT NULL THEN pass4
			--WHEN pass3 IS NOT NULL THEN pass3
			WHEN pass2 IS NOT NULL THEN pass2
			WHEN pass1 IS NOT NULL THEN pass1
		END - pass dur
		, pass1 - pass  dur1
		, pass2 - pass1 dur2
		, pass3 - pass2 dur3
		, pass4 - pass3 dur4
		, pass5 - pass4 dur5
	FROM Level5Passes 
), ResultsExt AS (
	SELECT r.laps, r.dur, r.dur1, r.dur2, r.dur3, r.dur4, r.dur5, p.number, p.name, p.category, p.team
	FROM Participants p
		LEFT JOIN Results r
			ON r.number = p.number
)
SELECT number
	, IFNULL(laps, 0) laps
	, IFNULL(dur , 0) dur
	, IFNULL(dur1, 0) dur1
	, IFNULL(dur2, 0) dur2
	, IFNULL(dur3, 0) dur3
	, IFNULL(dur4, 0) dur4
	, IFNULL(dur5, 0) dur5
	, IFNULL(name, '') name
	, IFNULL(category, '') category
	, IFNULL(team, '') team
FROM ResultsExt
ORDER BY CASE WHEN laps > 2 THEN 2 ELSE laps END DESC, dur`
)
