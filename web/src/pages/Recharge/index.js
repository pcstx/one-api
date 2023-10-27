import React, { useEffect, useState, useRef } from 'react';
import { Button, Form, Header, Segment, Statistic,Message,Confirm, FormGroup,Label,Image,Tab,Container } from 'semantic-ui-react';
import { API, showError, showWarning, showSuccess } from '../../helpers';
import { renderQuota } from '../../helpers/render'; 
import { getCookie, isH5,isMiniProgram,isWechat } from '../../helpers/utils'
import Script from 'react-load-script';

const Recharge = () => {
  const [redemptionCode, setRedemptionCode] = useState(1000);
  const [orderPrice,setOrderPrice] = useState(10);
  const [topUpLink, setTopUpLink] = useState('');
  const [userQuota, setUserQuota] = useState(0);
  const [userPoint, setUserPoint] = useState(0);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [open,setOpen] = useState(false);
  const [logo,setLogo] = useState('')
  const [hideImg,setHideImg] = useState(true)
 // const [query, setQuery] = useState(null);
 const intervalRef = useRef(); 

  let orderNumber=''
  let count = 199;

  const openConfirm = () =>  { 
    if(redemptionCode<100){
      showError('兑换最低100积分起')
      return;
    }
    if(redemptionCode>userPoint){
      showError('积分不足,请充值积分')
      return;
    }
    if(redemptionCode>1000000){
      showError('兑换最高100万积分')
      return;
    }
    setOpen(true) 
  }
  const closeConfirm = () => {setOpen(false)}

  const recharge = async () => {
    const redemptionCodeInt = parseInt(redemptionCode, 10);

    if(redemptionCodeInt<100){
      showError('兑换最低100积分起')
        return;
    }
    if(redemptionCodeInt>userPoint){
      showError('积分不足,请充值积分')
      return;
    }
    if(redemptionCodeInt>1000000){
      showError('兑换最高100万积分')
      return;
    }
    setIsSubmitting(true);
    try {
      const res = await API.post('/api/user/recharge', {
        point: redemptionCodeInt
      });
      const { success, message, data } = res.data;
      if (success) {
        showSuccess('积分兑换成功！');
        getUserQuota()
        // setUserQuota((quota) => {
        //   return quota + data;
        // });
        setRedemptionCode(100);
      } else {
        showError(message);
      }
    } catch (err) {
      showError('请求失败');
    } finally {
      setOpen(false)
      setIsSubmitting(false); 
    }
  };

  const cashRecharge = async () => {
    if(orderPrice<=0){
      showError('充值金额最低1元')
        return;
    }
    
    if(orderPrice>10000){
      showError('充值金额最高1万元')
      return;
    }

    setIsSubmitting(true);

    try {
      const _userAgent = window.navigator.userAgent
      if(isMiniProgram(_userAgent)){
        let user = localStorage.getItem('user');
        if (user) {
          let data = JSON.parse(user);
          window.wx.miniProgram.navigateTo({url: '/pages/redirect/redirect?orderPrice='+ orderPrice*1.0 +'&perkAIUserId='+ data.id})
        }
      } else if(isH5(_userAgent)){
        if (isWechat(_userAgent)) {
          const res = await API.post('/api/user/cashRecharge', {
            orderPrice: orderPrice*1.0,
            payType: 6
          });
          const { success, message, data } = res.data;
          if (success) {
            const { appId, nonceStr, paySign, prepayId, signType, timeStamp, _package } = data.payDataDto || {}
            window.wx.chooseWXPay({
                timestamp: timeStamp, // 支付签名时间戳，注意微信jssdk中的所有使用timestamp字段均为小写。但最新版的支付后台生成签名使用的timeStamp字段名需大写其中的S字符
                nonceStr, // 支付签名随机串，不长于 32 位
                package: _package, // 统一支付接口返回的prepay_id参数值，提交格式如：prepay_id=\*\*\*）
                signType, // 微信支付V3的传入RSA,微信支付V2的传入格式与V2统一下单的签名格式保持一致
                paySign, // 支付签名
                appId,
                success: function (r) {
                    //支付成功
                    if (r.errMsg == "chooseWXPay:ok") {
                      showSuccess('支付成功！');
                      getUserQuota()
                    } else {
                        //支付失败
                        showError('支付失败！');
                    }
                },
                fail: function (r) {
                    //支付失败
                    showError(r.errMsg ||'请求支付失败，稍后重新尝试！' )
                }
            })
          } else {
            showError(message);
          } 
        } else {
          const res = await API.post('/api/user/cashRecharge', {
            orderPrice: orderPrice*1.0,
            payType: 2
          });
          const { success, message, data } = res.data;
          if (success) {
            window.location.href = data.payUrl;             
          } else {
            showError(message);
          }
        }
      } else{
        const res = await API.post('/api/user/cashRecharge', {
          orderPrice: orderPrice*1.0,
          payType: 0
        });
        const { success, message, data } = res.data;
        if (success) {
          orderNumber = data.orderNumber
          setLogo(data.imgUrl)      
          setHideImg(false)
          //循环查询
          if(intervalRef.current){
            clearInterval(intervalRef.current)
          }
          const id  = setInterval(() => {
            queryOrder()
          }, 3000);
          intervalRef.current = id;
        } else {
          showError(message);
        }
      }
    } catch (err) {
      showError(err)
      showError('请求失败');
    } finally {
      setIsSubmitting(false); 
      setOpen(false);
    }
  };

  const openTopUpLink = () => {
    if (!topUpLink) {
      showError('超级管理员未设置充值链接！');
      return;
    }
    window.open(topUpLink, '_blank');
  };

  const getUserQuota = async ()=>{
    let res  = await API.get(`/api/user/point`);
    const {success, message, data} = res.data;
    if (success) {
      setUserQuota(data.quota);
      setUserPoint(data.point);
    } else {
      showError(message);
    }
  }

  const queryOrder = async () => {
    let res  = await API.get(`/api/user/queryOrder?orderNumber=${orderNumber}`);
    const {success, message, data} = res.data;
    if (success) { 
      if(data && data.orderStatus===1){
        //停止查询
        //增加金额
        clearInterval(intervalRef.current)
        showSuccess('充值成功！');
        setHideImg(true)
        setUserQuota((quota) => {
          return quota + data.orderPrice*500;
        });
      } else if(data && data.orderStatus === -1){
        //支付取消了
        clearInterval(intervalRef.current)
        showSuccess('支付已取消！');
        setHideImg(true)
      }
    } else {
      clearInterval(intervalRef.current)
      showError(message);
    }
    //循环超过次数
    if (count > 0) {
        count--;
    } else {    
        clearInterval(intervalRef.current);
        showError('支付超时！')
        setHideImg(true);
    }
  }

  const initWeChatSDK = async () => {
    let res = await API.get(`https://www.pushplus.plus/api/customer/wxApi/jsApi?url=https://perkai.pushplus.plus/recharge`,{
      headers: { "pushToken": getCookie('pushToken')}
    })
    const {code, msg, data} = res.data;
    if (code===200) { 
      if (window.wx && window.wx.config) {
          const { appId, nonceStr, signature, timestamp } = data
          window.wx.config({
              appId,
              nonceStr,
              signature,
              timestamp,
              jsApiList: ['chooseWXPay']
          })
      }
    } else{
      showWarning('微信SDK初始出错，稍后刷新页面重试~')
    }
  }
 
  useEffect(() => {
    initWeChatSDK()
    let status = localStorage.getItem('status');
    if (status) {
      status = JSON.parse(status);
      if (status.top_up_link) {
        setTopUpLink(status.top_up_link);
      }
    }
    getUserQuota().then(); 

    return () => { 
      //清空定时器
      if (intervalRef.current) {
        clearInterval(intervalRef.current)
      }
    };
  }, []);

  const Cash = () => {
    return (
      <>
        <Form>
        <FormGroup inline>
            <label>充值金额</label>
            <Form.Input 
              focus
              placeholder='元'
              name='orderPrice'
              value={orderPrice}
              onChange={(e) => {
                setOrderPrice(e.target.value);
              }}
              input='number'
              minLength='1'
              maxLength='5'
            />
            <Button color='orange' onClick={cashRecharge} disabled={isSubmitting}>
                {isSubmitting ? '充值中...' : '现金充值'}
            </Button>
           
        </FormGroup>
        <Label>
            充值Token数量：{(orderPrice*500*100).toLocaleString()}<br/>
        </Label>

        <Container className={`img-div  ${hideImg ? 'd-none':''}`}>
          <Image src={logo} size='small'/>
          <Label color='grey' className='label-div'>支持微信、支付宝扫码</Label>
        </Container>
       
    
    
            <Message info>
                <Message.List>
                    <Message.Item>现金1元=5万token</Message.Item>
                    <Message.Item>最少充值1元，最多1万元</Message.Item>
                    <Message.Item>50万token约为1美元额度</Message.Item>
                    <Message.Item>充值成功后无法退款</Message.Item>
                </Message.List>
            </Message>
        </Form>
        </>
    );
  }; 
   
  const Redemption = ()=> {
    return ( 
      <Form>
        <FormGroup inline>
            <label>积分数量</label>
            <Form.Input 
              focus
              placeholder='积分'
              name='redemptionCode'
              value={redemptionCode}
              onChange={(e) => {
                setRedemptionCode(e.target.value);
              }}
              input='number'
              minLength='3'
              maxLength='7'
            />
            <Button color='orange' onClick={openConfirm} disabled={isSubmitting}>
                {isSubmitting ? '兑换中...' : '兑换Token'}
            </Button>
          
        </FormGroup>
        <Label>
              兑换Token数量：{(redemptionCode*500).toLocaleString()}<br/>
            </Label>

            <Message info>
                <Message.List>
                    <Message.Item>积分兑换比例为1:500，100积分=5万token</Message.Item>
                    <Message.Item>最少100积分起兑，最多1百万积分</Message.Item>
                    <Message.Item>50万token约为1美元额度</Message.Item>
                    <Message.Item>兑换成token后将无法退回</Message.Item>
                </Message.List>
            </Message>
      </Form>        
    );
  };

  const panes = [
    { menuItem: '现金充值', render: () => 
    <Tab.Pane key='tab1'> 
      <Script
       url="//res.wx.qq.com/open/js/jweixin-1.6.0.js"
      />
    <Form>
    <FormGroup inline>
        <label>充值金额</label>
        <Form.Input 
          focus
          placeholder='元'
          name='orderPrice'
          value={orderPrice}
          onChange={(e) => {
            const inputValue = e.target.value;
            if (/^\d+(\.\d{0,2})?$/.test(inputValue)) {
              setOrderPrice(inputValue);
            }
          }}
          input='number'
          minLength='1'
          maxLength='5'
        />
        <Button color='orange' onClick={cashRecharge} disabled={isSubmitting}>
            {isSubmitting ? '充值中...' : '现金充值'}
        </Button>
       
    </FormGroup>
    <Label>
        充值Token数量：{(orderPrice*500*100).toLocaleString()}<br/>
    </Label>

    <Container className={`img-div  ${hideImg ? 'd-none':''}`}>
      <Image src={logo} size='small'/>
      <Label color='grey' className='label-div'>支持微信、支付宝扫码</Label>
    </Container>
   


        <Message info>
            <Message.List>
                <Message.Item>现金1元=5万token</Message.Item>
                <Message.Item>最少充值1元，最多1万元</Message.Item>
                <Message.Item>50万token约为1美元额度</Message.Item>
                <Message.Item>充值成功后无法退款</Message.Item>
            </Message.List>
        </Message>
    </Form></Tab.Pane>
  },
   { menuItem: '积分兑换', render: () => 
      <Tab.Pane key='tab2'><Form>
      <FormGroup inline>
          <label>积分数量</label>
          <Form.Input
            placeholder='积分'
            name='redemptionCode'
            value={redemptionCode}
            onChange={(e) => {
              const inputValue = e.target.value;
              if (/^[1-9]\d*$/.test(inputValue) || inputValue === '') {
                setRedemptionCode(inputValue);
              }
            }}
            input='number'
            minLength='3'
            maxLength='7'
          />
          <Button color='orange' onClick={openConfirm} disabled={isSubmitting}>
              {isSubmitting ? '兑换中...' : '兑换Token'}
          </Button>
        
      </FormGroup>
      <Label>
            兑换Token数量：{(redemptionCode*500).toLocaleString()}<br/>
          </Label>

          <Message info>
              <Message.List>
                  <Message.Item>积分兑换比例为1:500，100积分=5万token</Message.Item>
                  <Message.Item>最少100积分起兑，最多1百万积分</Message.Item>
                  <Message.Item>50万token约为1美元额度</Message.Item>
                  <Message.Item>兑换成token后将无法退回</Message.Item>
              </Message.List>
          </Message>
    </Form>    </Tab.Pane> 
    },
  ]

  return (
    <Segment>
      <Header as='h3'>额度充值</Header>

      <Statistic.Group widths='three'>
            <Statistic>
              <Statistic.Value>{userPoint.toLocaleString()}</Statistic.Value>
              <Statistic.Label>
                我的积分<br/>
                <a href='#' onClick={openTopUpLink}> 积分充值</a>
              </Statistic.Label>              
            </Statistic>
            <Statistic>
              <Statistic.Value>{userQuota.toLocaleString()}</Statistic.Value>
              <Statistic.Label>我的token</Statistic.Label>
            </Statistic>
            <Statistic>
              <Statistic.Value>{renderQuota(userQuota)}</Statistic.Value>
              <Statistic.Label>我的额度</Statistic.Label>
            </Statistic>
      </Statistic.Group>
 
      <Segment basic>
          <Tab menu={{ secondary: true, pointing: true }} panes={panes} />
      </Segment>

          <Confirm
            open={open}
            header='确认兑换'
            content='兑换后将无法退回，是否确认将积分兑换token？'
            onCancel={closeConfirm}
            onConfirm={recharge}
            cancelButton='取消'
            confirmButton="确认兑换"
            size='mini'
            />
    </Segment>
  );
};

export default Recharge;