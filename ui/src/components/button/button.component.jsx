import React from 'react';

import './button.styles.scss';

const Button = ({children}) => (
  <button className='button'>{children}</button>
);

export default Button;
