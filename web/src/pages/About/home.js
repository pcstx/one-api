import React, { useEffect, useState } from 'react';
import { Header, Segment,Container,Image } from 'semantic-ui-react'; 

const NewAbout = () => { 
  useEffect(() => {
    //displayAbout().then();
  }, []);

  return (
    <>
      {
         <Segment>
         <Header as='h2' textAlign='center'>使用说明</Header>
         <Header as={'h3'}>快速开始</Header>
         <Container>
          <p>1. 微信扫码登录系统。</p>
          <p>2. 点击顶部菜单栏中的“令牌”，点击“添加新的令牌”来创建自己的令牌。</p>
          <p>3. 在顶部菜单栏中点击“充值”，为自己账户中充值一定的token额度。</p>
          <p>4. 按照OpenAI相关接口将 https://api.openai.com 更换成 https://perkai.pushplus.plus ,再使用创建的令牌即可使用。</p>
          <p>5. 第三方系统可参考菜单栏的“聊天”，在“设置”中的“API Key”填入自己的令牌，“接口地址”中填入 https://perkai.pushplus.plus，即可正常使用。</p>
          <Image src={'/image/chat.png'} rounded bordered size='huge'/>
         </Container>
         
         <Header as={'h3'}>额度计算</Header>
         <Container>
          <p>请求接口将会消耗token，token的消耗按照字数来计算。可以在“日志"中查看具体消耗的额度。</p>
          <p>额度 =（提示 token 数 + 补全 token 数 * 补全倍率） * 模型倍率。其中补全倍率 GPT3.5 固定为 1.33，GPT4 为 2。</p>
          <p>调用非对话模型OpenAI接口会返回消耗的总 token，但是你要注意提示和补全的消耗倍率不一样。</p>
          <Image src={'/image/log.png'} rounded bordered size='huge'/>
          </Container>
       </Segment>
      }
    </>
  );
};


export default NewAbout;
