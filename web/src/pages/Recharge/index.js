import React, { useEffect, useState } from 'react';
import { Button, Form, Grid, Header, Segment, Statistic } from 'semantic-ui-react';
import { API, showError, showInfo, showSuccess } from '../../helpers';
import { renderQuota } from '../../helpers/render';

const Recharge = () => {
  const [redemptionCode, setRedemptionCode] = useState(100);
  const [topUpLink, setTopUpLink] = useState('');
  const [userQuota, setUserQuota] = useState(0);
  const [isSubmitting, setIsSubmitting] = useState(false);

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
    let res  = await API.get(`/api/user/self`);
    const {success, message, data} = res.data;
    if (success) {
      setUserQuota(data.quota);
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
      <Header as='h3'>积分兑换额度</Header>
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
            <Button color='green' onClick={openTopUpLink}>
              积分充值
            </Button>
            <Button color='yellow' onClick={recharge} disabled={isSubmitting}>
                {isSubmitting ? '兑换中...' : '兑换'}
            </Button>
          </Form>
        </Grid.Column>
        <Grid.Column>
          <Statistic.Group widths='one'>
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