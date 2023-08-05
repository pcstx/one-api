import React from 'react';

const Chat = () => {
  const chatLink = localStorage.getItem('chat_link')?localStorage.getItem('chat_link'):"https://chat.pushplus.plus/";

  return (
    <iframe
      src={chatLink}
      style={{ width: '100%', height: '85vh', border: 'none' }}
    />
  );
};


export default Chat;
