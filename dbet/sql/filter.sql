-- Prabha Dias’s DAS of Caulobacter
# SELECT t.tiltseriesID
# FROM TiltSeriesData AS t
#   JOIN UserData AS u ON u.DEF_id = t.`REF|UserData|user`
#   JOIN SpeciesData AS s ON s.DEF_id = t.`REF|SpeciesData|specie`
# WHERE u.fullname = 'Prabha Dias' AND s.SpeciesName LIKE 'Caulobacter%' AND t.ispublic = 1 limit 5 offset 5;



-- Gregory Henderson’s DAS of Ostreococcus tauri
# SELECT t.tiltseriesID
# FROM TiltSeriesData AS t
#   JOIN UserData AS u ON u.DEF_id = t.`REF|UserData|user`
#   JOIN SpeciesData AS s ON s.DEF_id = t.`REF|SpeciesData|specie`
# WHERE u.fullname = 'Gregory Henderson' AND s.SpeciesName = 'Ostreococcus tauri' AND t.ispublic = 1;



-- All public tomograms
# SELECT t.tiltseriesID
# FROM TiltSeriesData AS t
#   JOIN UserData AS u ON u.DEF_id = t.`REF|UserData|user`
#   JOIN SpeciesData AS s ON s.DEF_id = t.`REF|SpeciesData|specie`
# WHERE t.ispublic = 1 ORDER BY t.tiltseriesID ASC;



-- guoxi's fake
SELECT t.tiltseriesID
FROM TiltSeriesData AS t
  JOIN UserData AS u ON u.DEF_id = t.`REF|UserData|user`
  JOIN SpeciesData AS s ON s.DEF_id = t.`REF|SpeciesData|specie`
WHERE u.fullname = 'testuname' AND s.SpeciesName = 'testsname' AND t.ispublic = 1;