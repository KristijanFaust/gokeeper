import {gql} from '@apollo/react-hooks';

export default gql`
  mutation SignIn($email: String!, $password: String!) {
    signIn(input: {email:$email, password:$password}) {
      token
      user{
        id
        username
      }
    }
  }
`;