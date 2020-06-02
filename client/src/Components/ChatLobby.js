import React, {Component} from 'react';

class ChatLobby extends Component {
  constructor(props) {
    super(props);
    this.ws = new WebSocket("ws://localhost/ws") //CHANGE URL
    this.state = {
      message: "",
      messages: []
    };
  }

  renderChat = () => {
    return <p></p>
  }

  handleChange = (e) => {
    this.setState({message: e.target.value});
  }

  handleChat = (e) => {
    e.preventDefault();
    this.ws.send(this.state.message);
  }

  render() {
    

    return(
      <div>
        <div id="chat">{ this.ws.onmessage = (e) => { return <p>{e.data}</p> } }</div>
        <input type="text" id="chatInput" value={this.state.nickname} onChange={this.handleChange} />
        <button type="submit" id="chatButton" onClick={this.handleChat}>Chat</button>
      </div>
    );
  }
}

export default ChatLobby