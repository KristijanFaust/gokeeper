import {Link} from "react-router-dom";
import React from "react";

import './header-links.styles.scss'

const HeaderLinks = () => (
  <div className='header-link-group'>
    <Link className='header-link' to='/'>
      <span>Sign In</span>
    </Link>
    <div className='header-link-separator'>|</div>
    <Link className='header-link' to='/sign-up'>
      <span>Sign Up</span>
    </Link>
  </div>
);

export default HeaderLinks;