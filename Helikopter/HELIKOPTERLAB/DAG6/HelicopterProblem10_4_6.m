% TTK4135 - Helicopter lab
% Hints/template for problem 2.
% Updated spring 2017, Andreas L. Fl?ten

%% Initialization and model definition
init; % NB: Change this to the init file corresponding to your helicopter
delta_t	= 0.25; % sampling time

%Set as globals?
global N mx 


% Discrete time system model. x = [lambda r p p_dot e e_dot]'
A1 = eye(6)+ delta_t*[0 1 0 0 0 0;
        0 0 -K_2 0 0 0;
        0 0 0 1 0 0;
        0 0 -K_1*K_pp -K_1*K_pd 0 0;
        0 0 0 0 0 1;
        0 0 0 0 -K_3*K_ep -K_3*K_ed];
B1 = delta_t*[0 0; 0 0; 0 0; K_1*K_pp 0; 0 0; 0 K_3*K_ep];

% Number of states and inputs
mx = size(A1,2); % Number of states (number of columns in A)
mu = size(B1,2); % Number of inputs(number of columns in B)

% Initial values
x1_0 = pi;                              % Lambda
x2_0 = 0;                               % r
x3_0 = 0;                               % p
x4_0 = 0;                               % p_dot
x5_0 = 0;                               % e
x6_0 = 0;                               % e_dot
x0 = [x1_0 x2_0 x3_0 x4_0 x5_0 x6_0]';            % Initial values

% Time horizon and initialization
N  = 40;                               % Time horizon for states
M  = N;                                 % Time horizon for inputs
z  = zeros(N*mx+M*mu,1);                % Initialize z for the whole horizon
z0 = z;                                 % Initial value for optimization

% Bounds
ul 	    = -30*pi/180;                   % Lower bound on control -- u1
uu 	    = 30*pi/180;                    % Upper bound on control -- u1

xl      = -Inf*ones(mx,1);              % Lower bound on states (no bound)
xu      = Inf*ones(mx,1);               % Upper bound on states (no bound)
xl(3)   = ul;                           % Lower bound on state x3
xu(3)   = uu;                           % Upper bound on state x3

% Generate constraints on measurements and inputs
[vlb,vub]       = genbegr2(N,M,xl,xu,ul,uu);    % hint: genbegr2
vlb(N*mx+M*mu)  = 0;                            % We want the last input to be zero
vub(N*mx+M*mu)  = 0;                            % We want the last input to be zero

% Generate the matrix Q and the vector c (objecitve function weights in the QP problem) 
Q1 = zeros(mx,mx);
Q1(1,1) = 1;                             % Weight on state x1
Q1(2,2) = 0;                             % Weight on state x2
Q1(3,3) = 0;                             % Weight on state x3
Q1(4,4) = 0;                             % Weight on state x4
Q1(5,5) = 0;                             % Weight on state x5
Q1(6,6) = 0;                             % Weight on state x6
P1 = zeros(mu,mu);
P1(1,1) = 1;                             % Weight on input u1
P1(2,2) = 1;                             % Weight on input u2
Q = 2*genq2(Q1,P1,N,M,mu);              % Generate Q
%c = zeros(N*mx+M*mu,1);                 % Generate c

%% Generate system matrixes for linear model
Aeq = gena2(A1,B1,N,mx,mu);           % Generate A, hint: gena2
beq = zeros(N*mx,1);        	  % Generate b
beq(1:mx) = A1*x0; % Initial value

%% Solve SQP problem with nonlinear model

% alpha = 0.2;
% beta = 20;
% lambda_t = 2*pi/3;

options = optimoptions('fmincon','Algorithm','sqp');

f = @(z) z'*Q*z;

tic
z = fmincon(f,z0,[],[],Aeq,beq,vlb,vub,@constraint,options);
t1=toc;

% Calculate objective value
% phi1 = 0.0;
% PhiOut = zeros(N*mx+M*mu,1);
% for i=1:N*mx+M*mu
%   phi1=phi1+Q(i,i)*z(i)*z(i);
%   PhiOut(i) = phi1;
% end

%% Extract control inputs and states
u1  = [z(N*mx+1:mu:N*mx+M*mu-1)];
u2  = [z(N*mx+2:mu:N*mx+M*mu)];% Control input from solution

%u1  = [z(N*mx+1:mu:N*mx+M*mu);z(N*mx+M*mu-1)];
%u2  = [z(N*mx+2:mu:N*mx+M*mu);z(N*mx+M*mu)];% Control input from solution

x1 = [x0(1);z(mx+1:mx:N*mx)];              % State x1 from solution
x2 = [x0(2);z(mx+2:mx:N*mx)];              % State x2 from solution
x3 = [x0(3);z(mx+3:mx:N*mx)];              % State x3 from solution
x4 = [x0(4);z(mx+4:mx:N*mx)];              % State x4 from solution
x5 = [x0(5);z(mx+5:mx:N*mx)];              % State x5 from solution
x6 = [x0(6);z(mx+6:mx:N*mx)];              % State x6 from solution

num_variables = 5/delta_t;
zero_padding = zeros(num_variables,1);
unit_padding  = ones(num_variables,1);

u1   = [zero_padding; u1; zero_padding];
u2   = [zero_padding; u2; zero_padding];
x1  = [pi*unit_padding; x1; zero_padding];
x2  = [zero_padding; x2; zero_padding];
x3  = [zero_padding; x3; zero_padding];
x4  = [zero_padding; x4; zero_padding];
x5  = [zero_padding; x5; zero_padding];
x6  = [zero_padding; x6; zero_padding];

t = 0:delta_t:delta_t*(length(u1)-1);

% u1_andtime = [t' u1];
% u2_andtime = [t' u2];
u_star = [t' u1 u2];
x_star = [t' x1 x2 x3 x4 x5 x6];

%% LQR
 
Q_LQR = diag([1/pi^2 1 1/0.5^2 1/0.5^2 1/0.8^2 1/0.5^2]); %1 10 0.1 200 1 1
R_LQR = diag([1/0.53^2 1/0.37^2]);
[K,S,e] = dlqr(A1,B1,Q_LQR,R_LQR);



%% Plotting
%For t, see further up
t = 0:delta_t:delta_t*(length(u1)-1);


figure(2)
subplot(711)
stairs(t,u1),grid
ylabel('u')
subplot(712)
plot(t,x1,'m',t,x1,'mo'),grid
ylabel('lambda')
subplot(713)
plot(t,x2,'m',t,x2','mo'),grid
ylabel('r')
subplot(714)
plot(t,x3,'m',t,x3,'mo'),grid
ylabel('p')
subplot(715)
plot(t,x4,'m',t,x4','mo'),grid
ylabel('pdot')
subplot(716)
plot(t,x5,'m',t,x5,'mo'),grid
ylabel('e')
subplot(717)
plot(t,x6,'m',t,x6','mo'),grid
xlabel('tid (s)'),ylabel('edot')
