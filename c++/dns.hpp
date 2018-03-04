#ifndef DNS_HPP
#define DNS_HPP

#include <string>
#include "boost/asio.hpp"
using namespace boost;

namespace dns {

    struct Resolved {
        asio::ip::tcp::resolver::iterator ep_itr;
        boost::system::error_code ec;
    };

    template<typename T,
        typename = typename std::enable_if_t<
            std::is_constructible<std::string, T>::value
        >
    >
    Resolved Resolve(T&& domainName,
                     T&& portNum) noexcept {
        asio::io_service ios;

        // Creating a query
        asio::ip::tcp::resolver::query resolver_query(
            domainName, portNum,
            asio::ip::tcp::resolver::query::numeric_service
        );

        // Creating a resolver
        asio::ip::tcp::resolver resolver(ios);

        // used to store all information about resolution process
        // if resolution succeeded then it has endpoint iterator
        // otherwise, error_code
        dns::Resolved resolvedHost;

        // Now try to resolve domain name
        resolvedHost.ep_itr = resolver.resolve(
            resolver_query, resolvedHost.ec
        );

        return resolvedHost;
    }
}

#endif
