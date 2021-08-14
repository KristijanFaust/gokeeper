import {gql} from '@apollo/react-hooks';

export default gql`
  mutation SignUp($email: String!, $username:String!, $password: String!) {
    signUp(input: {email:$email, username: $username, password:$password}) {
      email
    }
  }
`;