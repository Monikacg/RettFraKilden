function [c,ceq] = constraint(z)
global N mx

alpha = 0.2;
beta = 20;
lambda_t = 2*pi/3;

c = zeros(N,1);

for k = 1:N
    lambda_k = z(mx*(k-1)+1);
    e_k = z(mx*(k-1)+5);
    c(k) = alpha*exp(-beta*(lambda_k-lambda_t).^2)-e_k;
end

ceq = [];
end 