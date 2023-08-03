import React, { useContext, useEffect, useState } from 'react';
import { Button, Divider, Form, Grid, Header, Image, Message, Modal, Segment } from 'semantic-ui-react';
import { Link, useNavigate, useSearchParams } from 'react-router-dom';
import { UserContext } from '../context/User';
import { API, getLogo, showError, showSuccess,setCookie } from '../helpers';
import Loading from '../components/Loading';

const QrCodeLogin = () => {
  const [inputs, setInputs] = useState({
    username: '',
    password: '',
    wechat_verification_code: ''
  });
  const [searchParams, setSearchParams] = useSearchParams();  
  const [qrCodeUrl, setQrCodeUrl] = useState({});
  const [qrCode, setQrCode] = useState({});
  const [userState, userDispatch] = useContext(UserContext);
  const [isDivVisible, setIsDivVisible] = useState(false);

  const logo = getLogo(); 
  let count = 19;
  let navigate = useNavigate();

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
                    setIsDivVisible(!isDivVisible);
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
        //console.log("data:"+ data);
        userDispatch({ type: 'login', payload: data });
        localStorage.setItem('user', JSON.stringify(data));
        navigate('/');
        showSuccess('登录成功！'); 

        return true;
      }else{
        //showError(message);
        return false;
      }
  }

  const handleToggleDiv = () => {
    // 切换div的显示状态
    getQrCodeUrl();
    count = 19;
    setIsDivVisible(!isDivVisible);    
  };
 
  return (
    <Grid textAlign='center' style={{ marginTop: '48px' }}>
      <Grid.Column style={{ maxWidth: 450 }}>
        <Header as='h2' color='' textAlign='center'>
          <Image src={logo} /> 用户登录
        </Header>
        <Form>
          <Segment>
          {qrCodeUrl && qrCodeUrl.length && qrCodeUrl.length>0 ?(<> <Image src={qrCodeUrl} size='medium' verticalAlign='middle'/></>):(<><Loading prompt={''}></Loading></>)}
            { isDivVisible && 
            <div  style={{backgroundColor: "rgba(255,255,255,.9)",
                width: "100%",
                height: "280px",
                top: "0",
                position: "absolute",
                top: "2rem",
                left: "0rem"
                }}>
                  <h2 style={{color: '#6c757d', height: "100%",padding: '130px 0 0 0',cursor: 'pointer'}} onClick={handleToggleDiv}>二维码过期，点击刷新</h2>
            </div>
            }
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
