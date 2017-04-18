function [c,ceq] = constraint(z)
N = 15;
mx = 6;

alpha = 0.2;
beta = 20;
lambda_t = 2*pi/3;

c = zeros(N,1);

for 1:N
    lambda_k = z(mx*(k-1)+1);
    e_k = z(mx*(k-1)+5);
    c(k) = alpha*exp(-beta*(lamba_k-lamba_t).^2)-e_k;
end

ceq = [];