import React, { useContext, useEffect, useState } from 'react';
import { Container,Button } from 'semantic-ui-react';
import { API, showError, showNotice } from '../../helpers';
import { useNavigate } from 'react-router-dom';
import { StatusContext } from '../../context/Status';

const NewHome = () => {
  const [statusState, statusDispatch] = useContext(StatusContext);
  const [homePageContentLoaded, setHomePageContentLoaded] = useState(false);
  const [homePageContent, setHomePageContent] = useState('');

  let navigate = useNavigate();

  const toAbout = () => navigate('/about');

  const toQQ = () => {
    const url = 'https://qm.qq.com/cgi-bin/qm/qr?k=1rVwgMbTO2JNsoSQnxZ_8mBn3W0zQbLI&jump_from=webapi&authKey=NUADwLOXuP8WY2J9qhzS8I7f5FNQw4vDDzDdW3l4vAz/wFYua1H6wvXHW+KYpmHS';
    window.open(url, '_blank');
  }

  const displayNotice = async () => {
    const res = await API.get('/api/notice');
    const { success, message, data } = res.data;
    if (success) {
      let oldNotice = localStorage.getItem('notice');
      if (data !== oldNotice && data !== '') {
        showNotice(data);
        localStorage.setItem('notice', data);
      }
    } else {
      showError(message);
    }
  };
 

  useEffect(() => {
    displayNotice().then(); 
  }, []);
  return (
    <>
      {
            <Container textAlign="center" style={{margin:"100px 0 100px 0",width:"755px"}}>
                <p style={{textAlign:'left',lineHeight:'30px'}}>
                    <strong style={{fontSize:'46px'}}>破壳AI</strong>
                </p>
                <p style={{textAlign:'left',lineHeight:'30px'}}>
                    <strong style={{fontSize:'46px'}}>生成式AI应用开发工具平台</strong>
                </p>
                <Container  style={{textAlign:'left',fontSize:'18px',lineHeight:'30px'}}>
                致力于整合AI能力，将其应用范围拓展到更多领域，使AI技术普及到更加广大人群。推动AI的全面发展，为社会带来积极的变革和发展机遇。
                </Container>
                <div>
                    <Button class="ui button" size="big" color="blue" onClick={toAbout} style={{margin:'60px 20px 60px 30px',backgroundColor:'#5680ff'}}>开始使用</Button>
                    <a href='https://qm.qq.com/cgi-bin/qm/qr?k=1rVwgMbTO2JNsoSQnxZ_8mBn3W0zQbLI&jump_from=webapi&authKey=NUADwLOXuP8WY2J9qhzS8I7f5FNQw4vDDzDdW3l4vAz/wFYua1H6wvXHW+KYpmHS'
                    style={{fontSize:'17px'}} target='_blank'>QQ群：161672147</a>
                </div>
            </Container>
      }

    </>
  );
};

export default NewHome;
