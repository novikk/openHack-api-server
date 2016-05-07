<!DOCTYPE html>
<html>
  <head>
    <style type="text/css">
      html, body { height: 100%; margin: 0; padding: 0; font-family: Helvetica}
      #map { height: calc(100% - 40px); }
    </style>
    <script src="https://code.jquery.com/jquery-2.2.3.min.js"></script>
  </head>
  <body>
    <div style="width: 100%; height: 20px; background: #333; color: white; padding: 10px; line-height: 20px">
      <span>EcoRun demo dashboard</span>
      <div style="float: right; padding-right: 30px; margin-top: -3px">
        <input type="text" name="km" style="padding: 5px;">
        <a href="#" onclick="newRoute()" style="margin-top: -1px; background: #A2E07B; padding: 7px; color: white; text-decoration: none; border-radius: 3px; border-bottom: 1px solid #17AD11; font-size: 12px">New route</a>
      </div>
    </div>
    <div id="map"></div>
    <script type="text/javascript">

      var map;
      var line;
      var marker;

      function initMap() {
        map = new google.maps.Map(document.getElementById('map'), {
          center: {lat: 43.46278, lng: -3.80500},
          zoom: 14
        });

        {{ range .data }}
        new google.maps.Circle({
          strokeColor: '#FF0000',
          strokeOpacity: 0.8,
          strokeWeight: 2,
          fillColor: '#FF0000',
          fillOpacity: 0.35,
          map: map,
          center: {lat: {{.Lat}}, lng: {{.Lng}}},
          radius: {{.Value}}
        });
        {{ end }}

        google.maps.event.addListener(map, 'click', function(event) {
           placeMarker(event.latLng);
        });
      }

      function placeMarker(location) {
          if (marker) marker.setMap(null);
          marker = new google.maps.Marker({
              position: location,
              map: map
          });
          console.log(location)
          url = '/route/new_with_end?tp=foot&km=0&fromLat=43.46278&fromLon=-3.80500&toLat=' +location.lat()+ "&toLon=" + location.lng()
          $.get( url, function( data ) {
            rt = JSON.parse(data.message);
            var ll = []
            for (var i = 0; i < rt.length; ++i) {
              ll.push({"lat": rt[i][0], "lng": rt[i][1]})
            }

            if (line) line.setMap(null);
            line = new google.maps.Polyline({
               path: ll,
               geodesic: true,
               strokeColor: '#FF0000',
               strokeOpacity: 1.0,
               strokeWeight: 2
             });

             line.setMap(map);
          });
      }

      function newRoute() {
        url = '/route/new?tp=foot' + "&km=" + Math.floor($("[name=km]").val()/2) + "&fromLat=43.46278" + "&fromLon=-3.80500"
        $.get( url, function( data ) {
          rt = JSON.parse(data.message);
          var ll = []
          for (var i = 0; i < rt.length; ++i) {
            ll.push({"lat": rt[i][0], "lng": rt[i][1]})
          }

          if (line) line.setMap(null);
          line = new google.maps.Polyline({
             path: ll,
             geodesic: true,
             strokeColor: '#FF0000',
             strokeOpacity: 1.0,
             strokeWeight: 2
           });

           line.setMap(map);
        });
      }

    </script>
    <script async defer
      src="https://maps.googleapis.com/maps/api/js?key=AIzaSyDb3m6o_OKV5OwUsLAdzdDQB8DQ6BJMIl0&callback=initMap">
    </script>
  </body>
</html>
