import React, { useEffect, useState } from 'react';
import { Button, Form, Grid, Header, Segment, Statistic,Message,Confirm } from 'semantic-ui-react';
import { API, showError, showInfo, showSuccess } from '../../helpers';
import { renderQuota } from '../../helpers/render';

const Recharge = () => {
  const [redemptionCode, setRedemptionCode] = useState(100);
  const [topUpLink, setTopUpLink] = useState('');
  const [userQuota, setUserQuota] = useState(0);
  const [userPoint, setUserPoint] = useState(0);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [open,setOpen] = useState(false);

  const openConfirm = () =>  { setOpen(true) }
  const closeConfirm = () => {setOpen(false)}

  const recharge = async () => {
    if(redemptionCode<100){
        showInfo('兑换最低100积分起')
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
      <Header as='h3'>积分兑换token</Header>
      <Grid columns={2} stackable>
        <Grid.Column>
          <Form>
            <Form.Input
              placeholder='积分'
              name='redemptionCode'
              value={redemptionCode}
              onChange={(e) => {
                setRedemptionCode(e.target.value);
              }}
            />
            <Message info>
                <Message.List>
                    <Message.Item>积分兑换比例为1:500，100积分=50,000 token</Message.Item>
                    <Message.Item>最少100积分起兑</Message.Item>
                    <Message.Item>500,000 token约为1美元</Message.Item>
                    <Message.Item>兑换成token后将无法退回！</Message.Item>
                </Message.List>
            </Message>
            <Button color='orange' onClick={openConfirm} disabled={isSubmitting}>
                {isSubmitting ? '兑换中...' : '兑换额度'}
            </Button>
            <Button color='blue' onClick={openTopUpLink}>
              积分充值
            </Button>
          </Form>
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
        </Grid.Column>
        <Grid.Column>
          <Statistic.Group widths='three'>
            <Statistic>
              <Statistic.Value>{userPoint.toLocaleString()}</Statistic.Value>
              <Statistic.Label>我的积分</Statistic.Label>
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
        </Grid.Column>
      </Grid>
    </Segment>
  );
};

export default Recharge;