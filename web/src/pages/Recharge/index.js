import React, { useEffect, useState } from 'react';
import { Button, Form, Grid, Header, Segment, Statistic,Message,Confirm, FormGroup,Container,Label } from 'semantic-ui-react';
import { API, showError, showInfo, showSuccess } from '../../helpers';
import { renderQuota } from '../../helpers/render';

const Recharge = () => {
  const [redemptionCode, setRedemptionCode] = useState(100);
  const [topUpLink, setTopUpLink] = useState('');
  const [userQuota, setUserQuota] = useState(0);
  const [userPoint, setUserPoint] = useState(0);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [open,setOpen] = useState(false);

  const openConfirm = () =>  { 
    if(redemptionCode<100){
      showError('兑换最低100积分起')
      return;
    }
    if(redemptionCode<userPoint){
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
    if(redemptionCode<userPoint){
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
      setIsSubmitting(false); 
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

  useEffect(() => {
    let status = localStorage.getItem('status');
    if (status) {
      status = JSON.parse(status);
      if (status.top_up_link) {
        setTopUpLink(status.top_up_link);
      }
    }
    getUserQuota().then();
  }, []);

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
                      <Message.Item>积分兑换比例为1:500，100积分=50,000 token</Message.Item>
                      <Message.Item>最少100积分起兑，最多1,000,000积分</Message.Item>
                      <Message.Item>500,000 token约为1美元额度</Message.Item>
                      <Message.Item>兑换成token后将无法退回</Message.Item>
                  </Message.List>
              </Message>
        </Form>
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