import SignIn from '../../components/sign-in/sign-in.component';

import './sign-in-page.styles.scss';

const SignInPage = ({signInCallback}) => (
  <div className='sign-in-page'>
    <SignIn signInCallback={signInCallback} />
  </div>
);

export default SignInPage;
