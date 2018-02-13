#include <iostream>
#include "dns.hpp"

int main(int argc, char** argv) {
    auto resolved(dns::Resolve(
        std::string(argv[1]), std::string("3333")
    ));
    
    if (resolved.ec != 0) {
        std::cerr << "Failed to resolve a DNS name."
                  << "\nError code: " << resolved.ec.value()
                  << "\nMessage: " << resolved.ec.message() << std::endl;
    }

    asio::ip::tcp::resolver::iterator ep_end;

    for (; resolved.ep_itr != ep_end; ++resolved.ep_itr) {
        auto endpoint(resolved.ep_itr->endpoint());

        std::cout << endpoint << std::endl;
    }

    return 0;
}
