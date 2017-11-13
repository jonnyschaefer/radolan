## Notice
*Data Source:* ftp://ftp-cdc.dwd.de/pub/CDC/grids_germany/hourly/radolan/RADOLAN-Koordinatendateien%20Lambda%20Phi/

```
RADOLAN-ME-Raster.zip/E_products_lambda_bottom.txt        -> ./lambda_bottom_1500x1400.txt 
Raster-lambda_phi_left_bottom.zip/lambda_bottom.txt       -> ./lambda_bottom_900x900.txt   
lambda-phi_center_1100x900.zip/lambda_center_1100x900.txt -> ./lambda_center_1100x900.txt  
RADOLAN-ME-Raster.zip/E_products_lambda_center.txt        -> ./lambda_center_1500x1400.txt 
Raster-lambda_phi_center.zip/lambda_center.txt            -> ./lambda_center_900x900.txt   
RADOLAN-ME-Raster.zip/E_products_phi_bottom.txt           -> ./phi_bottom_1500x1400.txt    
Raster-lambda_phi_left_bottom.zip/phi_bottom.txt          -> ./phi_bottom_900x900.txt      
lambda-phi_center_1100x900.zip/phi_center_1100x900.txt    -> ./phi_center_1100x900.txt     
RADOLAN-ME-Raster.zip/E_products_phi_center.txt           -> ./phi_center_1500x1400.txt    
Raster-lambda_phi_center.zip/phi_center.txt               -> ./phi_center_900x900.txt      
```

## Begin of original readme
Readme zu den Dateien lambda_bottom.txt und phi_bottom.txt bzw. lambda_center.txt und phi_center.txt

Version 1.1, 10.12.2004

Daniel Sacher

MeteoSolutions GmbH
Sturzstr. 45
64285 Darmstadt

Tel.: 06151-599 03 42


Die Datei lambda_bottom.txt / phi_bottom.txt  enthält die geografische Länge / geografische Breite
der linken unteren Ecke jedes Pixels im Deutschland-Komposit.

Die Datei lambda_center.txt / phi_center.txt  enthält die geografische Länge / geografische Breite
des Zentralpunktes jedes Pixels im Deutschland-Komposit.

 
Die räumliche Auflösung beträgt 1km x 1km.

Die Daten sind zeilenweise abgespeichert (900 Werte in 900 Zeilen), beginnend mit 
der linken unteren Ecke, also kopfüber.
Das Dezimalformat ist: F8.5 (FORTRAN-Bezeichner)

Zur Berechnung wurden folgende, innerhalb des DWD übliche Formeln verwendet:

	inverse Polarstereografische Projektion

        LON=arctan(-x/y)+LON0

	LAT= 2arctan(R A sin(LON-LON0)/x)-90°

        mit R: Erdradius R=6370,04 km
            A: A=1+sin(LAT0)
            LON0: 10° E
            LAT0: 60° N

        kartesische Koordinaten der Pixel (Koordinatenursprung in (LON0, LAT0)):

        x=R M cos(LAT)sin(LON-LON0)

        y=-R M cos(LAT)cos(LON-LON0)
         
        mit 
            M=(1+sin(LAT0))/(1+sin(LAT))
        

