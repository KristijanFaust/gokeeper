import {Link} from 'react-router-dom';

import {ReactComponent as Emblem} from '../../../assets/logo.svg';

import './logo.styles.scss';

const Logo = () => (
  <Link className='logo' to='/'>
    <Emblem className='emblem'/>
    <h1>GoKeeper</h1>
  </Link>
);

export default Logo;