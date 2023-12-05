import React from 'react';

const Application = () => {
  return (
    <div class="ui cards"> 
    <div class="card">
      <div class="content">
        <img class="right floated mini ui image" src="/icon.png"/>
        <div class="header">聊天</div>
        <div class="meta">chatAI</div>
        <div class="description">多模型智能聊天机器人</div>
      </div>
      <a href='https://chat.pushplus.plus/' target='_blank'>
        <div class="ui bottom attached button"><i class="add icon"></i> 点击使用 </div>
      </a>
    </div>

    <div class="card">
      <div class="content">
        <img class="right floated mini ui image" src="/icon.png"/>
        <div class="header">文生图</div>
        <div class="meta">createImage</div>
        <div class="description">通过简单的文字描述就能生成你想要的图片</div>
      </div>
      <a href='http://perkai.pushplus.plus/app/tti' target='_blank'>
        <div class="ui bottom attached button"><i class="add icon"></i> 点击使用 </div>
      </a>
    </div>

    <div class="card">
      <div class="content">
        <img class="right floated mini ui image" src="/icon.png"/>
        <div class="header">网页生成器</div>
        <div class="meta">HTML Page Generator</div>
        <div class="description">根据您绘制的线框生成 html 页面</div>
      </div>
      <a href='http://html.pushplus.plus' target='_blank'>
        <div class="ui bottom attached button"><i class="add icon"></i> 点击使用 </div>
      </a>
    </div>

    {/* <div class="card">
      <div class="content">
        <img class="right floated mini ui image" src="/icon.png"/>
        <div class="header">知识库</div>
        <div class="meta">基于向量数据库方案</div>
        <div class="description">基于大语言模型的AI知识库问答平台</div>
      </div>
      <a href='https://knowledge.pushplus.plus' target='_blank'>
        <div class="ui bottom attached button"><i class="add icon"></i> 点击使用 </div>
      </a>
    </div>
 
    <div class="card">
      <div class="content">
        <img class="right floated mini ui image" src="/icon.png"/>
        <div class="header">图生图</div>
        <div class="meta">
        </div>
        <div class="description">基于现有图片生成同类图片</div>
      </div>
      <a href='https://knowledge.pushplus.plus' target='_blank'>
        <div class="ui bottom attached button"><i class="add icon"></i> 点击使用 </div>
      </a>
    </div>



    <div class="card">
      <div class="content">
        <img class="right floated mini ui image" src="/icon.png"/>
        <div class="header">Logo生成器</div>
        <div class="meta">
        </div>
        <div class="description">简单绘制几笔草图智能生成logo</div>
      </div>
      <a href='https://knowledge.pushplus.plus' target='_blank'>
        <div class="ui bottom attached button"><i class="add icon"></i> 点击使用 </div>
      </a>
    </div> */}

    <div class="card">
      <div class="content">
        <img class="right floated mini ui image" src="/icon.png"/>
        <div class="header">pushplus 推送加</div>
        <div class="meta">微信公众号能力</div>
        <div class="description">集成微信、短信、邮件、Bark、腾讯轻联、钉钉、飞书等渠道的信息推送平台</div>
      </div>
      <a href='https://www.pushplus.plus' target='_blank'>
        <div class="ui bottom attached button"><i class="add icon"></i> 点击使用 </div>
      </a>
    </div>
 
    <div class="card">
      <div class="content">
        <img class="right floated mini ui image" src="/icon.png"/>
        <div class="header">微加机器人</div>
        <div class="meta">wechat</div>
        <div class="description">微信消息辅助机器人,可以通过接口主动的让微信下发消息。</div>
      </div>
      <a href='https://bot.pushplus.plus' target='_blank'>
        <div class="ui bottom attached button"><i class="add icon"></i> 点击使用 </div>
      </a>
    </div>
 
  </div>
  );
};


export default Application;
