var results = [];

// var scoreList = React.render(
//   <ScoreList results={results} />,
//   document.getElementById('example')
// );

var sock = new SockJS('/sock');
sock.onmessage = function(e) {
  var msg = JSON.parse(e.data);
  switch(msg.type){
    case 'init':
      results = msg.data.slice();
      break;
    case 'diff':
      
      console.log(msg.data);
      
      // results = window.applyDiff(results, msg.data);
      var txtArea = document.getElementById("responseText");
      txtArea.value = JSON.stringify(msg.data);
      
      var cnt = document.getElementById("cnt");
      cnt.value = msg.data.length;
      
      break;
  }
  // scoreList.setProps({ results: results });
};

function incrementScore(id) {
  sock.send(JSON.stringify({
    type: 'increment',
    id: id
  }));
}

function deletePlayer(id) {
  sock.send(JSON.stringify({
    type: 'delete',
    id: id
  }));
}

function insertPlayer(name) {
  sock.send(JSON.stringify({
    type: 'insert',
    name: name
  }));
}
