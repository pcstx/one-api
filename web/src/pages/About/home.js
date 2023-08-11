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
          <p>4. 按照OpenAI相关接口将 https://api.openai.com 更换成 <strong>https://perkai.pushplus.plus</strong> ,再使用创建的令牌即可使用。</p>
          <p>5. 第三方系统可参考菜单栏的“聊天”，在“设置”中的“令牌”填入自己的令牌，即可正常使用。如有“接口地址”或“baseUrl”则填入 <strong>https://perkai.pushplus.plus</strong>。</p>
          <Image src={'https://image.pushplus.plus/image/chat.png'} rounded bordered size='huge'/>
         </Container>
         
         <Header as={'h3'}>额度计算</Header>
         <Container>
          <p>请求接口将会消耗token，token的消耗按照字数来计算。可以在“日志"中查看具体消耗的额度。</p>
          <p>额度 =（提示 token 数 + 补全 token 数 * 补全倍率） * 模型倍率。其中补全倍率 GPT3.5 固定为 1.33，GPT4 为 2。</p>
          <p>调用非对话模型OpenAI接口会返回消耗的总 token，但是你要注意提示和补全的消耗倍率不一样。</p>
         
          <p>预估费用如下：</p>
          <table class='ui small basic compact table'>
            <thead>
                <tr>
                    <th>大语言</th>
                    <th>模型(Model)</th>
                    <th>输入文字(Input)</th>
                    <th>输出文字(Output)</th>
                </tr>
            </thead>
            <tbody>
                <tr>
                    <td>GPT-3.5 Turbo</td>
                    <td>4K context</td>
                    <td>	$0.0015 / 1K tokens</td>
                    <td>	$0.002 / 1K tokens</td>
                </tr>
                <tr>
                    <td>GPT-3.5 Turbo</td>
                    <td>16K context</td>
                    <td>	$0.003 / 1K tokens</td>
                    <td>	$0.004 / 1K tokens</td>
                </tr>
                <tr>
                    <td>GPT-4</td>
                    <td>8K context	</td>
                    <td>	$0.03 / 1K tokens</td>
                    <td>	$0.06 / 1K tokens</td>
                </tr>
                <tr>
                    <td>GPT-4</td>
                    <td>32K context</td>
                    <td>	$0.06 / 1K tokens</td>
                    <td>	$0.12 / 1K tokens</td>
                </tr>
            </tbody>
        </table>

        <Image src={'/image/log.png'} rounded bordered size='huge'/>
          </Container>
       </Segment>
      }
    </>
  );
};


export default NewAbout;
