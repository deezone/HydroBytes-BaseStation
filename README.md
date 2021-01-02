# HydroBytes-BaseStation
The Base Station is a part of a collection of systems called
"[HydroBytes](https://github.com/deezone/HydroBytes)" that uses micro
controllers to manage and monitor plant health in an automated garden.

## Introduction

The "garden" is simply a backyard patio in Brooklyn, New York. Typically
there are only two seasons - cold and hot in New York City. By
automating an urban garden ideally the space will thrive with minimum
supervision. The amount of effort to automate is besides the point,  
everyone needs their vices.

- **[Water Management Station](https://github.com/deezone/HydroBytes-WaterManagement)**
- **Base Station**
- **[Plant Station](https://github.com/deezone/HydroBytes-PlantStation)**

![brooklyn-20201115 garden layout](https://raw.githubusercontent.com/deezone/HydroBytes/master/resources/gardenBrooklynDiagram-20201115.jpg)

![Garden](https://github.com/deezone/HydroBytes-WaterManagement/blob/master/resources/garden-01.png)

### YouTube Channel

[![YouTube Channel](https://github.com/deezone/HydroBytes-WaterManagement/blob/master/resources/youTube-TN.png?raw=true)](https://www.youtube.com/channel/UC00A_lEJD2Hcy9bw6UuoUBA "All of the HydroBytes videos")

### Notes

Development of a Go based API is based on instruction in the amazing
courses at **[Ardan Labs](https://education.ardanlabs.com/collections?category=courses)**.

#### Starting Web Server
```
go run .
2021/01/01 19:08:42 Listening on localhost:8000
```

- request to `localhost:8000`:
![Basic GET response](https://github.com/deezone/HydroBytes-BaseStation/blob/master/resources/images/basic-GET-response.jpg?raw=true)
