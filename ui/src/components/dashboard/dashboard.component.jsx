import {useState} from 'react';
import {useQuery} from '@apollo/react-hooks';

import Spinner from '../spinner/spinner.component';
import ErrorMessage from '../messages/error/error-message.component';
import PasswordEditor from './password-editor/password-editor.component';
import CreatePasswordForm from './create-password-form/create-password-form.component';
import userPasswordsQuery from '../../graphql/queries/query-user-passwords';

import './dashboard.component.scss';

const Dashboard = () => {
  const userId = localStorage.getItem('userId');
  const [passwords, setPasswords] = useState([]);
  const [errors, setErrors] = useState(null);

  const {loading} = useQuery(userPasswordsQuery, {
    variables: {userId: userId},
    onCompleted: (data) => {
      setPasswords(data.queryUserPasswords);
    },
    onError: (response) => {
      setErrors(response.graphQLErrors.map(error => error.message));
      if (!response.graphQLErrors.length) {
        setErrors(["Server is unavailable!"])
      }
    }
  })

  const addNewPassword = (newPassword) => {
    setPasswords([...passwords, newPassword]);
  };

  if (loading) return <div className='dashboard'><Spinner /></div>;

  const passwordsTitle = (passwords.length || errors != null) ? <span className='category-header'>Passwords</span> :
    <span className='category-header'>Seems like you haven't saved any passwords yet</span>;

  const passwordsEditor = passwords.length ? <PasswordEditor passwords={passwords}/> : null;

  const errorMessage = errors ? errors.map((error, index) => {
    return <ErrorMessage key={index}>{error}</ErrorMessage>
  }) : null;

  return (
    <div className='dashboard'>
      {passwordsTitle}
      {errorMessage}
      {passwordsEditor}
      <span className='category-header'>Add new password</span>
      <CreatePasswordForm createPasswordsCallback={addNewPassword}/>
    </div>
  );
};

export default Dashboard;