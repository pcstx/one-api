import React, { useEffect, useState } from 'react';
import { Button, Form, Header, Segment, Statistic,Message,Confirm, FormGroup,Label,Image,Tab,Container } from 'semantic-ui-react';
import { API, showError, showSuccess } from '../../helpers';
import { renderQuota } from '../../helpers/render'; 

const Recharge = () => {
  const [redemptionCode, setRedemptionCode] = useState(100);
  const [orderPrice,setOrderPrice] = useState();
  const [topUpLink, setTopUpLink] = useState('');
  const [userQuota, setUserQuota] = useState(0);
  const [userPoint, setUserPoint] = useState(0);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [open,setOpen] = useState(false);
  const [logo,setLogo] = useState('')
  const [hideImg,setHideImg] = useState(true)
  const [query, setQuery] = useState(null);
  
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
    setIsSubmitting(true);
    try {
      const res = await API.post('/api/user/recharge', {
        point: redemptionCode
      });
      const { success, message, data } = res.data;
      if (success) {
        showSuccess('积分兑换成功！');
        setUserQuota((quota) => {
          return quota + data;
        });
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
      const res = await API.post('/api/user/cashRecharge', {
        orderPrice: orderPrice*1.0
      });
      const { success, message, data } = res.data;
      if (success) {
        orderNumber = data.orderNumber
        setLogo(data.imgUrl)      
        setHideImg(false)
        //循环查询
        if(query){
          clearInterval(query)
        }
        setQuery(setInterval(() => {
          queryOrder()
        }, 3000))
      } else {
        showError(message);
      }
    } catch (err) {
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
        clearInterval(query)
        showSuccess('充值成功！');
        setHideImg(true)
        setUserQuota((quota) => {
          return quota + data.orderPrice*500;
        });
      } else if(data && data.orderStatus === -1){
        //支付取消了
        clearInterval(query)
        showSuccess('支付已取消！');
        setHideImg(true)
      }
    } else {
      showError(message);
    }
    //循环超过次数
    if (count > 0) {
        count--;
    } else {    
        clearInterval(query);
        showSuccess('支付超时！');
        setHideImg(true);
    }
  }
 
  useEffect(() => {
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
      if (query) {
        clearInterval(query)
      }
    };
  }, [query]);

  const Cash = () => {
    return (
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
    <Tab.Pane key='tab1'> <Form>
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