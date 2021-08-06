import React from "react";

import Logo from '../logo/logo.component';
import HeaderLinks from "./header-links/header-links.component";
import UserPanel from "./user-panel/user-panel.component";

import './header.styles.scss';

const Header = () => (
    <div className='header'>
      <Logo />
      <HeaderLinks />
      <UserPanel />
    </div>
);

export default Header;