import React from 'react';

import Input from '../input/input.component';
import Button from '../button/button.component';

import './sign-in.styles.scss';

const SignIn = () => (
  <div className='sign-in'>
    <h2>Sign in</h2>

    <form>
      <Input name='email' type='email' label='email' required/>
      <Input name='password' type='password' label='password' required/>
      <Button type='submit'> Sign in </Button>
    </form>
  </div>
);

export default SignIn;
