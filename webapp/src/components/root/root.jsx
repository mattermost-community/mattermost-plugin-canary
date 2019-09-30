import React from 'react';
import PropTypes from 'prop-types';

export default class Root extends React.Component {
    static propTypes = {
        checkCanaryCookie: PropTypes.func.isRequired,
    }
    componentDidMount() {
        this.props.checkCanaryCookie();
    }
    render() {
        return null;
    }
}