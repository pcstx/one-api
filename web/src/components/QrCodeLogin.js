import React, { useContext, useEffect, useState } from 'react';
import { Button, Divider, Form, Grid, Header, Image, Message, Modal, Segment } from 'semantic-ui-react';
import { Link, useNavigate, useSearchParams } from 'react-router-dom';
import { UserContext } from '../context/User';
import { API, getLogo, showError, showSuccess } from '../helpers';

const QrCodeLogin = () => {
  const [inputs, setInputs] = useState({
    username: '',
    password: '',
    wechat_verification_code: ''
  });
  const [searchParams, setSearchParams] = useSearchParams();  
  const [qrCodeUrl, setQrCodeUrl] = useState({});
  const [qrCode, setQrCode] = useState({});
  const logo = getLogo(); 
  let count = 5;

  useEffect(() => {
    if (searchParams.get('expired')) {
      showError('未登录或登录已过期，请重新登录！');
    }
    let status = localStorage.getItem('status');
    if (status) {
      status = JSON.parse(status);
    }
    getQrCodeUrl();
  }, []); 
 
  async function getQrCodeUrl() {
    const res = await API.get(
      `/api/user/wechat/getQrcode`
    );
    const { success, message, data } = res.data;
    if (success) {
      setQrCodeUrl(data.qrCodeUrl);
      
        //轮询确认登录
        let timer = setInterval(async () => {
            let isSuccess = await confirmLogin(data.qrCode); 
            if(isSuccess){
                clearInterval(timer);
            } else{
                if (count > 0) {
                    count--;
                } else {    
                    clearInterval(timer);
                }
            }
        }, 3000);
      return data;  
    }
  };

  async function confirmLogin(qrCode){
    const res = await API.get(
        `/api/user/wechat/confirmLogin?key=${qrCode}`
      );
      const { success, message, data } = res.data;
      if (success) {
        console.log(data);
        //在cookie中创建key为pushToken，值为data内容
        document.cookie = "pushToken=" + data + ";path=/";
        return true;
      }
      
      return false;
  }
 
  return (
    <Grid textAlign='center' style={{ marginTop: '48px' }}>
      <Grid.Column style={{ maxWidth: 450 }}>
        <Header as='h2' color='' textAlign='center'>
          <Image src={logo} /> 用户登录
        </Header>
        <Form size='large'>
          <Segment>
            <Image src={qrCodeUrl} size='medium' />
            <div style={{ textAlign: 'center' }}>
                <h3>
                  微信扫码登录
                </h3>
              </div> 
          </Segment>
        </Form> 
      </Grid.Column>
    </Grid>
  );
};

export default QrCodeLogin;
